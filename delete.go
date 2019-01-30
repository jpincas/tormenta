package tormenta

import (
	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

func (db DB) Delete(entity Record, ids ...gouuidv6.UUID) (int, error) {
	// If a separate entity ID has been specified then use it
	if len(ids) > 0 {
		entity.SetID(ids[0])
	}

	// Deleted counter
	var deleted int

	err := db.KV.Update(func(txn *badger.Txn) error {

		root := KeyRoot(entity)
		key := newContentKey(root, entity.GetID()).bytes()

		// First check to see if the key actually exists
		_, err := txn.Get(key)
		if err != badger.ErrKeyNotFound {
			if err := txn.Delete(key); err != nil {
				return err
			}

			deleted++
		}

		return nil
	})

	return deleted, err
}
