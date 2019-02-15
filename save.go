package tormenta

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dgraph-io/badger"
)

const (
	errNoModel = "Cannot save entity %s - it does not have a tormenta model"
)

func (db DB) Save(entities ...Record) (int, error) {

	err := db.KV.Update(func(txn *badger.Txn) error {
		for _, entity := range entities {
			// Make a copy of the entity and attempt to get the old
			// version from the DB for deindexing
			newEntity := newRecord(entity)
			found, err := db.Get(newEntity, entity.GetID())
			if err != nil {
				return err
			}

			// If it does exist, then we'll need to deindex it.
			// If it's a new entity then deindexing is not necessary
			if found {
				if err := deIndex(txn, newEntity); err != nil {
					return err
				}
			}

			// Presave trigger
			if err := entity.PreSave(); err != nil {
				return err
			}

			// Build the key root
			keyRoot, e := entityTypeAndValue(entity)

			// Check that the model field exists
			modelField := e.FieldByName("Model")
			if !modelField.IsValid() {
				return fmt.Errorf(errNoModel, keyRoot)
			}

			// Assert the model type
			// Check if there is an idea, if not create one
			// Update the time last updated
			model := modelField.Interface().(Model)
			if model.ID.IsNil() {
				model.ID = newID()
			}
			model.LastUpdated = time.Now().UTC()

			// Set the new model back on the entity
			modelField.Set(reflect.ValueOf(model))

			data, err := db.serialise(entity)
			if err != nil {
				return err
			}

			key := newContentKey(keyRoot, model.ID).bytes()
			if err := txn.Set(key, data); err != nil {
				return err
			}

			// Post save trigger
			entity.PostSave()

			// indexing
			if err := index(txn, entity); err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return len(entities), nil
}

// The regular 'Save' function is atomic - if there is any error, the whole thing
// gets rolled back.  If you don't care about atomicity, you can use SaveIndividually

// SaveIndividually discards atomicity and continues saving entities even if one fails.
// The total count of saved entities is returned.
// Badger transactions have a maximum size, so the regular 'Save' function is best used
// for a small number of entities.  This function could be used to save 1 million entities
// if required
func (db DB) SaveIndividually(entities ...Record) (counter int, lastErr error) {
	for _, entity := range entities {
		if _, err := db.Save(entity); err != nil {
			lastErr = err
		} else {
			counter++
		}
	}

	return counter, lastErr
}
