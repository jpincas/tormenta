package tormenta

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dgraph-io/badger"
)

type tormentable interface {
	MarshalMsg([]byte) ([]byte, error)
	UnmarshalMsg([]byte) ([]byte, error)
}

const (
	errNoModel = "Cannot save entity %s - it does not have a tormenta model"
)

func (db DB) Save(entities ...tormentable) (int, error) {
	err := db.KV.Update(func(txn *badger.Txn) error {
		for _, entity := range entities {
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

	return len(entities), nil
}
