package tormentadb

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dgraph-io/badger"
)

const (
	errNoModel = "Cannot save entity %s - it does not have a tormenta model"
)

func (db DB) Save(entities ...Tormentable) (int, error) {

	noBatches := noBatches(len(entities), batchSize)
	for i := 0; i < noBatches; i++ {

		err := db.KV.Update(func(txn *badger.Txn) error {
			a, b := batchStartAndEnd(i, batchSize, len(entities))
			batch := entities[a:b]

			for _, entity := range batch {
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

				entityMsg, err := entity.MarshalMsg(nil)
				if err != nil {
					return err
				}

				key := newContentKey(keyRoot, model.ID).bytes()
				if err := txn.Set(key, entityMsg); err != nil {
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

	}

	return len(entities), nil
}
