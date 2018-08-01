package tormentadb

import (
	"errors"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

func (db DB) Delete(root string, ids ...gouuidv6.UUID) (int, error) {
	if len(ids) == 0 {
		return 0, errors.New("Must specify at least one ID to delete")
	}

	// Deleted counter
	var deleted int

	err := db.KV.Update(func(txn *badger.Txn) error {
		for _, id := range ids {

			key := newContentKey([]byte(root), id).bytes()

			// First check to see if the key actually exists
			_, err := txn.Get(key)
			if err != badger.ErrKeyNotFound {
				if err := txn.Delete(key); err != nil {
					return err
				}
				deleted++
			}
		}

		return nil
	})

	return deleted, err
}
