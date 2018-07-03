package tormenta

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dgraph-io/badger"
)

type Tormentable interface {
	MarshalMsg([]byte) ([]byte, error)
	UnmarshalMsg([]byte) ([]byte, error)
}

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
				// Build the key root
				keyRoot, e := getKeyRoot(entity)

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

				entityMsg, _ := entity.MarshalMsg(nil)
				txn.Set(makeKey(keyRoot, model.ID), entityMsg)
			}

			return nil
		})

		if err != nil {
			return 0, err
		}

	}

	return len(entities), nil
}
