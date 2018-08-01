package tormentadb

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

const (
	ErrNoID = "Cannot get entity %s - ID is nil"
)

// Get retrieves an entity, either according to the ID set on the entity,
// or using a separately specified ID (optional)
func (db DB) Get(entity Tormentable, ids ...gouuidv6.UUID) (bool, error) {
	keyRoot, e := entityTypeAndValue(entity)

	// Check that the model field exists
	modelField := e.FieldByName("Model")
	if !modelField.IsValid() {
		return false, fmt.Errorf(errNoModel, keyRoot)
	}

	// Assert the model
	model := modelField.Interface().(Model)

	// If an ID has been specified, overwrite the one on the entity
	var id gouuidv6.UUID
	if len(ids) > 0 {
		id = ids[0]
	} else {
		if model.ID.IsNil() {
			return false, fmt.Errorf(ErrNoID, keyRoot)
		}
		id = model.ID
	}

	err := db.KV.View(func(txn *badger.Txn) error {
		key := newContentKey(keyRoot, id).bytes()
		item, err := txn.Get(key)

		if err != nil {
			return err
		}

		val, err := item.Value()
		if err != nil {
			return err
		}

		_, err = entity.UnmarshalMsg(val)
		if err != nil {
			return err
		}

		// Post Get trigger
		entity.PostGet()

		return nil
	})

	if err == badger.ErrKeyNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// Post Get trigger and set Created field
	entity.GetCreated()
	entity.PostGet()

	return true, nil
}
