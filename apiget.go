package tormenta

import (
	"github.com/jpincas/gouuidv6"
)

var noCTX = make(map[string]interface{})

// Get retrieves an entity, either according to the ID set on the entity,
// or using a separately specified ID (optional, takes priority)
func (db DB) Get(entity Record, ids ...gouuidv6.UUID) (bool, error) {
	return db.GetWithContext(entity, noCTX, ids...)
}

// GetWithContext retrieves an entity, either according to the ID set on the entity,
// or using a separately specified ID (optional, takes priority), and allows the passing of a non-empty context.
func (db DB) GetWithContext(entity Record, ctx map[string]interface{}, ids ...gouuidv6.UUID) (bool, error) {
	txn := db.KV.NewTransaction(false)
	defer txn.Discard()

	return db.get(txn, entity, ctx, ids...)
}

func (db DB) GetIDs(target interface{}, ids ...gouuidv6.UUID) (int, error) {
	return db.GetIDsWithContext(target, noCTX, ids...)
}

func (db DB) GetIDsWithContext(target interface{}, ctx map[string]interface{}, ids ...gouuidv6.UUID) (int, error) {
	txn := db.KV.NewTransaction(false)
	defer txn.Discard()

	return db.getIDsWithContext(txn, target, ctx, ids...)
}
