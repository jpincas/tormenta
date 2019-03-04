package tormenta

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/dgraph-io/badger"

	"github.com/jpincas/gouuidv6"
)

const (
	ErrIDFieldNotExist             = "%s field was not found"
	ErrFieldNotExist               = "%s field was not found"
	ErrFieldWrongType              = "%s field should be either a pointer to a struct, or slice thereof"
	ErrIDFieldIncorrectType        = "%s is not an ID field of the type UUID"
	ErrNoRecords                   = "at least 1 record is needed in order to load relations"
	ErrRelationMustBeStructPointer = "relation must be a pointer to a struct"

	idFieldPostfixSingle   = "ID"
	idFieldPostfixMultiple = "IDs"
	fieldPathSep           = "."
)

// TODO: this code is currently horribly complex in terms of
// reflection and concurrency and needs to be refactored.
// But first we need a proper test suite on all the relational stuff

func idFieldName(fieldName string, postfix string) string {
	return fieldName + postfix
}

func fieldPath(fieldName string) []string {
	return strings.Split(fieldName, fieldPathSep)
}

// reJoinFieldPath is a bit counterintutive
// When we recursively call HasOne for nested relations,
// we need to join back up the field path (minus the 1st component)
// which has already been dealt with.  We then call HasOne with the rejoined
// path as the SINGLE member of the 'relationsToLoad' argument
func reJoinFieldPath(pathComponents []string) []string {
	return []string{strings.Join(pathComponents, fieldPathSep)}
}

// QueryModifier is a function that modifies and existing query
type QueryModifier = func(q *Query) *Query

func LoadByQuery(db *DB, fieldName string, queryModifier QueryModifier, entities ...Record) error {
	if len(entities) == 0 {
		return fmt.Errorf("LoadByQuery requires at least one entity")
	}

	txn := db.KV.NewTransaction(false)
	defer txn.Discard()

	type relationsQueryResult struct {
		entityID    gouuidv6.UUID
		relationIDs idList
		err         error
	}

	ch := make(chan relationsQueryResult)
	defer close(ch)
	var wg sync.WaitGroup

	// Do some reflect work up front on the first entity
	exampleEntity := entities[0]

	// Reflect on the specified field - bail if its not there
	field := recordValue(exampleEntity).FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf(ErrFieldNotExist, fieldName)
	}

	indexString := indexStringForThisEntity(exampleEntity) + idFieldPostfixSingle

	// Now we can go ahead and get results for all the IDs
	// The related field on the entity is a slice of pointers,
	// but queries expect a pointer to a slice,
	// So we do a bit of reflect magic to get a target for the query
	target := reflect.New(reflect.SliceOf(field.Type().Elem().Elem())).Interface()

	for i := range entities {
		wg.Add(1)
		go func(ii int) {
			query := db.Find(target).Match(indexString, entities[ii].GetID())

			// If a query modifier has been passed in, use it to add a new AND clause
			if queryModifier != nil {
				query = queryModifier(query)
			}

			ids, err := query.queryIDs(txn)
			ch <- relationsQueryResult{
				entityID:    entities[ii].GetID(),
				relationIDs: ids,
				err:         err,
			}
		}(i)
	}

	relatedEntitiesForEachEntity := map[gouuidv6.UUID]idList{}
	allIDsToGet := idList{}
	var errorsList []error

	go func() {
		for relationsQueryResult := range ch {
			if relationsQueryResult.err != nil {
				errorsList = append(errorsList, relationsQueryResult.err)
			} else {
				relatedEntitiesForEachEntity[relationsQueryResult.entityID] = relationsQueryResult.relationIDs
				allIDsToGet = append(allIDsToGet, relationsQueryResult.relationIDs...)
			}

			wg.Done()
		}
	}()

	wg.Wait()

	if len(errorsList) > 0 {
		return errorsList[0]
	}

	db.getIDsWithContext(txn, target, noCTX, allIDsToGet...)

	// Once we have all the results,
	// we build up a map of results keyed by ID
	resultsMap := map[gouuidv6.UUID]Record{}
	s := reflect.ValueOf(target).Elem()
	for i := 0; i < s.Len(); i++ {
		e := s.Index(i).Addr().Interface().(Record)
		resultsMap[e.GetID()] = e
	}

	// Now, finally, for each entity we originally passed in
	// we can go back, retrieve the list of its related entities
	// from the first map that resulted from the query,
	// and then for each of those IDs, get the actual result
	// from the final results map and append it the target slice on the entity
	// (as a slice of pointers)
	for i := range entities {
		thisEntityRelatedIDs := relatedEntitiesForEachEntity[entities[i].GetID()]

		// Up until now we have been taking the 'field' from the first entity passed in as an exmple
		// Now we actually have to operate on that field for each entity
		field := recordValue(entities[i]).FieldByName(fieldName)

		for _, relatedID := range thisEntityRelatedIDs {
			field.Set(
				reflect.Append(
					field,
					reflect.ValueOf(resultsMap[relatedID]),
				),
			)

		}
	}

	return nil
}

type relationsResult struct {
	fieldName string
	recordMap map[gouuidv6.UUID]Record
	err       error
}

func LoadByID(db *DB, relationsToLoad []string, entities ...Record) error {
	txn := db.KV.NewTransaction(false)
	defer txn.Discard()
	return loadByID(db, txn, relationsToLoad, entities...)
}

func loadByID(db *DB, txn *badger.Txn, relationsToLoad []string, entities ...Record) error {
	// We need at least 1 entity to make this work
	if len(entities) == 0 {
		return errors.New(ErrNoRecords)
	}

	ch := make(chan relationsResult)
	defer close(ch)

	var wg sync.WaitGroup

	// For each fieldname/path specified for relational loading,
	// we spawn a worker to go and get all the relations needed
	// for ALL the entities - we'll do the sorting and reattaching later
	for _, relation := range relationsToLoad {
		path := fieldPath(relation)

		if len(path) == 0 {
			return nil
		}

		wg.Add(1)
		go func(thisPath []string) {
			recordMap, err := getRelatedField(db, txn, thisPath[0], entities...)

			// If there is more than one component to the path,
			// call HasOne recursively, passing in the rest of the path
			// (joined back up with the separator, and passed a single
			// member of a slice)
			// and the entities that came back above
			if len(thisPath) > 1 {
				var nestedEntities []Record
				for _, record := range recordMap {
					nestedEntities = append(nestedEntities, record)
				}

				if err := loadByID(db, txn, reJoinFieldPath(path[1:]), nestedEntities...); err != nil {
					log.Println("error in nested HasOne")
					// TODO: need to work out way of signaling this at top level
				}
			}

			// Wait until the nested loading has finished
			// before sending the result to the channel,
			// otherwise the top level loading will finish before the lower level
			ch <- relationsResult{
				fieldName: thisPath[0],
				recordMap: recordMap,
				err:       err,
			}
		}(path)
	}

	// The workers return a map of relational records keyed by ID,
	// As the results come back, we'll build up a 'master' map
	// of those relation maps, keyed by the field name
	masterRecordMap := map[string]map[gouuidv6.UUID]Record{}
	var errorsList []error
	go func() {
		for relationsResult := range ch {
			if relationsResult.err != nil {
				errorsList = append(errorsList, relationsResult.err)
			} else {
				masterRecordMap[relationsResult.fieldName] = relationsResult.recordMap
			}

			wg.Done()
		}
	}()

	// Once all the relations are in, bail if there was any errorr
	wg.Wait()
	if len(errorsList) > 0 {
		return errorsList[0]
	}

	// At this point we have a 'master' map that contains all the relations
	// we need for each field requested and for all the entities.
	// Now we have to go through each entity, and for each field requested, retrieve
	// that record and 'attach' it according to the stored xxxID or xxxIDs field.
	// We do that in parallel for each entity
	var entityWg sync.WaitGroup
	entityWg.Add(len(entities))

	done := make(chan bool)
	defer close(done)

	for i := range entities {
		go func(ii int) {
			for fieldName, recordMap := range masterRecordMap {
				// First thing is to work out whether this is a single or multiple relation get
				resultsField := recordValue(entities[ii]).FieldByName(fieldName)

				switch resultsField.Type().Kind() {
				case reflect.Ptr:
					field := recordValue(entities[ii]).FieldByName(idFieldName(fieldName, idFieldPostfixSingle))

					// No need to confirm that the interface to UUID is OK
					// as this is performed already in the inner loop so will
					// always be OK at this point
					id := field.Interface().(gouuidv6.UUID)

					// Check that the record is in the map and set if it is
					// Since the inner loop won't add zero identifiers to the map,
					// this also acts as a sanity check against those zero identifiers
					record, relationFound := recordMap[id]
					if relationFound {
						recordValue(entities[ii]).FieldByName(fieldName).Set(reflect.ValueOf(record))
					}

				case reflect.Slice:
					// For slices, things are slightly more complex,
					// as we need to iterate all the related IDs and append the results one by one to the
					// results target slice
					field := recordValue(entities[ii]).FieldByName(idFieldName(fieldName, idFieldPostfixMultiple))
					ids := field.Interface().([]gouuidv6.UUID)
					for _, id := range ids {
						record, relationFound := recordMap[id]
						if relationFound {
							// Nasty code - a reflect append!!
							recordValue(entities[ii]).FieldByName(fieldName).Set(
								reflect.
									Append(recordValue(entities[ii]).
										FieldByName(fieldName), reflect.ValueOf(record)),
							)
						}
					}
				}
			}

			done <- true
		}(i)
	}

	go func() {
		for range done {
			entityWg.Done()
		}
	}()

	entityWg.Wait()
	return nil
}

func getRelatedField(db *DB, txn *badger.Txn, fieldName string, entities ...Record) (map[gouuidv6.UUID]Record, error) {

	recordMap := map[gouuidv6.UUID]Record{}

	// For each entity, add the ID of the relation to the map
	// giving deduping for free
	for _, entity := range entities {

		// Work out whether this is a single ID or a list of IDs
		// by inspecting the resultsField name - bail if its not there
		resultsField := recordValue(entity).FieldByName(fieldName)
		if !resultsField.IsValid() {
			return recordMap, fmt.Errorf(ErrFieldNotExist, fieldName)
		}

		switch resultsField.Type().Kind() {
		case reflect.Ptr:
			idfieldName := idFieldName(fieldName, idFieldPostfixSingle)
			field := recordValue(entity).FieldByName(idfieldName)
			if !field.IsValid() {
				return recordMap, fmt.Errorf(ErrIDFieldNotExist, idfieldName)
			}

			id, ok := field.Interface().(gouuidv6.UUID)
			if !ok {
				return recordMap, fmt.Errorf(ErrIDFieldIncorrectType, idfieldName)
			}

			if !id.IsNil() {
				recordMap[id] = nil
			}

		case reflect.Slice:
			idfieldName := idFieldName(fieldName, idFieldPostfixMultiple)
			field := recordValue(entity).FieldByName(idfieldName)
			if !field.IsValid() {
				return recordMap, fmt.Errorf(ErrIDFieldNotExist, idfieldName)
			}

			ids, ok := field.Interface().([]gouuidv6.UUID)
			if !ok {
				return recordMap, fmt.Errorf(ErrIDFieldIncorrectType, idfieldName)
			}

			for _, id := range ids {
				if !id.IsNil() {
					recordMap[id] = nil
				}
			}

		default:
			return recordMap, fmt.Errorf(ErrFieldWrongType, fieldName)
		}
	}

	// id map -> list
	ids := make([]gouuidv6.UUID, 0, len(recordMap))
	for k := range recordMap {
		ids = append(ids, k)
	}

	// Now we have the IDs of the related entities we need to get,
	// we just have to work out what type we are getting.
	// Use the first record as an exemplar -
	fieldValue := fieldValue(entities[0], fieldName)

	// Again, how we proceed here depends on whether this is a single or multiple
	// relation get
	var typeToGet reflect.Type

	switch fieldValue.Type().Kind() {
	case reflect.Ptr:
		typeToGet = fieldValue.Type().Elem()

	case reflect.Slice:
		// double Elem(): slice -> pointer -> actual type
		typeToGet = fieldValue.Type().Elem().Elem()
	}

	// Set up a new slice of the type we are getting
	// and use the multiple Get by ID api to grab all the
	// relations

	results := newSlice(typeToGet, len(ids))
	if _, err := db.getIDsWithContext(txn, results, noCTX, ids...); err != nil {
		return recordMap, err
	}

	// At this point, results is *[]WhateverTheEntityIs
	// We'll iterate it and return as a map of *WhateverTheEntityIs
	// which fulfuls the Record interface
	s := reflect.ValueOf(results).Elem()
	for i := 0; i < s.Len(); i++ {
		record := s.Index(i).Addr().Interface().(Record)
		recordMap[record.GetID()] = record
	}

	return recordMap, nil
}
