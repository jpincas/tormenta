package tormenta

import (
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

const (
	ErrRecordNotFound = "Record with ID %v was not found"
)

func (db DB) Delete(entity Record, ids ...gouuidv6.UUID) error {
	// If a separate entity ID has been specified then use it
	if len(ids) > 0 {
		entity.SetID(ids[0])
	}

	// First lets try to get the entity,
	// Its a good sanity check to make sure it really exists,
	// but more importantly we're going to need to deindex it,
	// so we'll need it current state
	if found, err := db.Get(entity); err != nil {
		return err
	} else if !found {
		return fmt.Errorf(ErrRecordNotFound, entity.GetID())
	}

	err := db.KV.Update(func(txn *badger.Txn) error {
		if err := deleteRecord(txn, entity); err != nil {
			return err
		}

		// if err := deIndexRecord(txn, entity); err != nil {
		// 	return err
		// }

		return nil
	})

	return err
}

func deleteRecord(txn *badger.Txn, entity Record) error {
	root := KeyRoot(entity)
	key := newContentKey(root, entity.GetID()).bytes()
	return txn.Delete(key)
}
