package tormenta

import (
	"errors"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

func (db DB) Delete(root string, ids ...gouuidv6.UUID) (int, error) {
	if len(ids) == 0 {
		return 0, errors.New("Must specify at least one ID to delete")
	}

	err := db.KV.Update(func(txn *badger.Txn) error {
		for _, id := range ids {
			if err := txn.Delete(newContentKey([]byte(root), id).bytes()); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return len(ids), nil
}
