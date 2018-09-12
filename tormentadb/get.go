package tormentadb

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

const (
	ErrNoID = "Cannot get entity %s - ID is nil"
)

// Get retrieves an entity, either according to the ID set on the entity,
// or using a separately specified ID (optional)
func (db DB) Get(entity Tormentable, ids ...gouuidv6.UUID) (bool, int, error) {
	return db.get(entity, make(map[string]interface{}), ids...)
}

// Get retrieves an entity, either according to the ID set on the entity,
// or using a separately specified ID (optional), this is the same as the above 'Get'
// function but allows the passing of a non-empty context.
func (db DB) GetWithContext(entity Tormentable, ctx map[string]interface{}, ids ...gouuidv6.UUID) (bool, int, error) {
	return db.get(entity, ctx, ids...)
}


func (db DB) get(entity Tormentable, ctx map[string]interface{}, ids ...gouuidv6.UUID) (bool, int, error) {
	// start the timer
	t0 := time.Now()

	keyRoot, e := entityTypeAndValue(entity)

	// Check that the model field exists
	modelField := e.FieldByName("Model")
	if !modelField.IsValid() {
		return false, timerMiliseconds(t0), fmt.Errorf(errNoModel, keyRoot)
	}

	// Assert the model
	model := modelField.Interface().(Model)

	// If an ID has been specified, overwrite the one on the entity
	var id gouuidv6.UUID
	if len(ids) > 0 {
		id = ids[0]
	} else {
		if model.ID.IsNil() {
			return false, timerMiliseconds(t0), fmt.Errorf(ErrNoID, keyRoot)
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
		// This seems to be a duplicate action
		// so commenting for now
		// entity.PostGet(ctx)

		return nil
	})

	if err == badger.ErrKeyNotFound {
		return false, timerMiliseconds(t0), nil
	} else if err != nil {
		return false, timerMiliseconds(t0), err
	}

	// Populate the created field
	entity.GetCreated()

	// Run the 'postGet' trigger
	entity.PostGet(ctx)

	return true, timerMiliseconds(t0), nil
}
