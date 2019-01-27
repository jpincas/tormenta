package tormenta

import (
	"reflect"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

const (
	ErrNoID       = "Cannot get entity %s - ID is nil"
	ErrTooManyIDs = "You should only specify 1 ID for a get request"
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
	// Get the key root and cache the value
	_, value := entityTypeAndValue(target)

	// Set up list results placeholder
	results := reflect.Indirect(reflect.ValueOf(target))

	var counter int
	for _, id := range ids {
		var entity Record
		entity = reflect.New(value.Type().Elem()).Interface().(Record)

		// For an error, we'll bail,
		// but if we simply can't find the record, we'll continue
		if found, err := db.get(entity, noCTX, id); err != nil {
			return counter, err
		} else if found {
			results = reflect.Append(
				results,
				reflect.Indirect(reflect.ValueOf(entity)),
			)
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
