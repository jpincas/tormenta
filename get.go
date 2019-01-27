package tormenta

import (
	"reflect"

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

func (db DB) GetIDs(target interface{}, ids ...gouuidv6.UUID) (int, error) {
	records := newResultsArray(target)

	var counter int
	for _, id := range ids {
		// Its inefficient creating a new entity target for the result
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
