package tormenta

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	errNoID = "Cannot get entity %s - ID is nil"
)

func (db DB) Get(entity tormentable) (bool, error) {
	keyRoot, e := getKeyRoot(entity)

	// Check that the model field exists
	modelField := e.FieldByName("Model")
	if !modelField.IsValid() {
		return false, fmt.Errorf(errNoModel, keyRoot)
	}

	// Assert the model
	model := modelField.Interface().(Model)
	if model.ID.IsNil() {
		return false, fmt.Errorf(errNoID, keyRoot)
	}

	err := db.KV.View(func(txn *badger.Txn) error {
		item, err := txn.Get(makeKey(keyRoot, model.ID))

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

		return nil
	})

	if err == badger.ErrKeyNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil

}
