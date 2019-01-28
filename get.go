package tormenta

import (
	"reflect"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

const (
	ErrNoID = "Cannot get entity %s - ID is nil"
)

var noCTX = make(map[string]interface{})

// Get retrieves an entity, either according to the ID set on the entity,
// or using a separately specified ID (optional, takes priority)
func (db DB) Get(entity Record, ids ...gouuidv6.UUID) (bool, error) {
	return db.get(entity, noCTX, ids...)
}

// GetWithContext retrieves an entity, either according to the ID set on the entity,
// or using a separately specified ID (optional, takes priority), and allows the passing of a non-empty context.
func (db DB) GetWithContext(entity Record, ctx map[string]interface{}, ids ...gouuidv6.UUID) (bool, error) {
	return db.get(entity, ctx, ids...)
}

type getResult struct {
	id     gouuidv6.UUID
	record Record
	found  bool
	err    error
}

func (db DB) GetIDs(target interface{}, ctx map[string]interface{}, ids ...gouuidv6.UUID) (int, error) {
	return db.GetIDsWithContext(target, noCTX, ids...)
}

func (db DB) GetIDsWithContext(target interface{}, ctx map[string]interface{}, ids ...gouuidv6.UUID) (int, error) {
	ch := make(chan getResult)
	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)

		// It's inefficient creating a new entity target for the result
		// on every loop, but we can't just create a single one
		// and reuse it, because there would be risk of data from 'previous'
		// entities 'infecting' later ones if a certain field wasn't present
		// in that later entity, but was in the previous one.
		// Unlikely if the all JSON is saved with the schema, but I don't
		// think we can risk it
		go func(thisRecord Record, thisID gouuidv6.UUID) {
			found, err := db.get(thisRecord, ctx, thisID)
			ch <- getResult{
				id:     thisID,
				record: thisRecord,
				found:  found,
				err:    err,
			}
		}(newRecord(target), id)
	}

	var resultsList []Record
	var errorsList []error
	go func() {
		for getResult := range ch {
			if getResult.err != nil {
				errorsList = append(errorsList, getResult.err)
			} else if getResult.found {
				resultsList = append(resultsList, getResult.record)
			}

			// Only signal to the wait group that a record has been fetched
			// at this point rather than the anonymous func above, otherwise
			// you tend to lose the last result
			wg.Done()
		}
	}()

	// Once all the results are in, we need to
	// sort them according to the original order
	// But we'll bail now if there were any errors
	wg.Wait()

	if len(errorsList) > 0 {
		return 0, errorsList[0]
	}

	return sortToOriginalIDsOrder(target, resultsList, ids), nil
}

func sortToOriginalIDsOrder(target interface{}, resultList []Record, ids []gouuidv6.UUID) (counter int) {
	resultMap := map[gouuidv6.UUID]Record{}
	for _, record := range resultList {
		resultMap[record.GetID()] = record
	}

	records := newResultsArray(target)

	// Remember, we didn't bail if a record was not found
	// so there is a chance it won't be in the map - thats ok - just keep count
	// of the ones that are there
	for _, id := range ids {
		record, found := resultMap[id]
		if found {
			records = reflect.Append(records, recordValue(record))
			counter++
		}
	}

	return counter
}

func (db DB) get(entity Record, ctx map[string]interface{}, ids ...gouuidv6.UUID) (bool, error) {
	// If an override id has been specified, set it on the entity
	if len(ids) > 0 {
		entity.SetID(ids[0])
	}

	err := db.KV.View(func(txn *badger.Txn) error {
		item, err := txn.Get(newContentKey(KeyRoot(entity), entity.GetID()).bytes())
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) {
			// TODO: unmarshalling error?
			db.json.Unmarshal(val, entity)
		})
	})

	// We are not treating 'not found' as an actual error,
	// instead we return 'false' and nil (unless there is an actual error)
	if err == badger.ErrKeyNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// Get triggers
	entity.GetCreated()
	entity.PostGet(ctx)

	return true, nil
}

// For benchmarking / comparison with parallel get

func (db DB) GetIDsSerial(target interface{}, ids ...gouuidv6.UUID) (int, error) {
	records := newResultsArray(target)

	var counter int
	for _, id := range ids {
		// It's inefficient creating a new entity target for the result
		// on every loop, but we can't just create a single one
		// and reuse it, because there would be risk of data from 'previous'
		// entities 'infecting' later ones if a certain field wasn't present
		// in that later entity, but was in the previous one.
		// Unlikely if the all JSON is saved with the schema, but I don't
		// think we can risk it
		record := newRecord(target)

		// For an error, we'll bail, if we simply can't find the record, we'll continue
		if found, err := db.get(record, noCTX, id); err != nil {
			return counter, err
		} else if found {
			records = reflect.Append(records, recordValue(record))
			counter++
		}
	}

	return counter, nil
}
