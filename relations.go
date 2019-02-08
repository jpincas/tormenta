package tormenta

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/jpincas/gouuidv6"
)

const (
	ErrIDFieldNotExist             = "%s field was not found"
	ErrIDFieldIncorrectType        = "%s is not an ID field of the type UUID"
	ErrNoRecords                   = "at least 1 record is needed in order to load relations"
	ErrRelationMustBeStructPointer = "relation must be a pointer to a struct"

	idFieldPostFix = "ID"
	fieldPathSep   = "."
)

// TODO: clean this up when all relational stuff is done

func idFieldName(fieldName string) string {
	return fieldName + idFieldPostFix
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

type relationsResult struct {
	fieldName string
	recordMap map[gouuidv6.UUID]Record
	err       error
}

// HasOne
func HasOne(db *DB, relationsToLoad []string, entities ...Record) error {

	fmt.Println(relationsToLoad)

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
			recordMap, err := hasOne(db, thisPath[0], entities...)

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

				if err := HasOne(db, reJoinFieldPath(path[1:]), nestedEntities...); err != nil {
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
	// that record and 'attach' it according to the stored xxxID field.
	// We do that in parallel for each entity
	var entityWg sync.WaitGroup
	entityWg.Add(len(entities))

	done := make(chan bool)
	defer close(done)

	for i := range entities {
		go func(ii int) {
			for fieldName, recordMap := range masterRecordMap {
				field := recordValue(entities[ii]).FieldByName(idFieldName(fieldName))
				// No need to confirm that the interface to UUID is OK
				// as this is performed already in the inner loop so will
				// always be OK at this point
				id := field.Interface().(gouuidv6.UUID)

				// Get the record from the record map - if its nil
				// don't worry, the relation will just be nil
				// fmt.Println("output: ", ToJSON(entities[ii]))
				recordValue(entities[ii]).FieldByName(fieldName).Set(reflect.ValueOf(recordMap[id]))
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

func hasOne(db *DB, fieldName string, entities ...Record) (map[gouuidv6.UUID]Record, error) {
	idfieldName := idFieldName(fieldName)
	recordMap := map[gouuidv6.UUID]Record{}

	// For each entity, add the ID of the relation to the map
	// giving deduping for free
	for _, entity := range entities {
		field := recordValue(entity).FieldByName(idfieldName)
		if !field.IsValid() {
			return recordMap, fmt.Errorf(ErrIDFieldNotExist, idfieldName)
		}

		id, ok := field.Interface().(gouuidv6.UUID)
		if !ok {
			return recordMap, fmt.Errorf(ErrIDFieldIncorrectType, idfieldName)
		}

		recordMap[id] = nil
	}

	// id map -> list
	ids := make([]gouuidv6.UUID, 0, len(recordMap))
	for k := range recordMap {
		ids = append(ids, k)
	}

	// Now we have the IDs of the related entities we need to get,
	// we just have to work out what type we are getting.
	// Use the first record as an exemplar -
	// check that its a struct
	fieldValue := fieldValue(entities[0], fieldName)

	if fieldValue.Kind() != reflect.Ptr {
		return recordMap, errors.New(ErrRelationMustBeStructPointer)
	}

	if reflect.ValueOf(fieldValue).Kind() != reflect.Struct {
		return recordMap, errors.New(ErrRelationMustBeStructPointer)
	}

	// Set up a new slice of the type we are getting
	// and use the multiple Get by ID api to grab all the
	// relations

	results := newSlice(fieldValue.Type().Elem(), len(ids))
	if _, err := db.GetIDs(results, ids...); err != nil {
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
