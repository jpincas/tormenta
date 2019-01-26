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
		// a, b := batchStartAndEnd(i, batchSize, len(entities))
		// batch := entities[a:b]

		for _, entity := range entities {
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

			data, err := json.Marshal(entity)
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
			if err := index(txn, entity, keyRoot, model.ID); err != nil {
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
