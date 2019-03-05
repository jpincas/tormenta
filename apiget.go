package tormenta

import (
	"time"

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
	t := time.Now()

	txn := db.KV.NewTransaction(false)
	defer txn.Discard()

	ok, err := db.get(txn, entity, ctx, ids...)

	if db.Options.DebugMode {
		var n int
		if ok {
			n = 1
		}
		debugLogGet(entity, t, n, err, ids...)
	}

	return ok, err
}

func (db DB) GetIDs(target interface{}, ids ...gouuidv6.UUID) (int, error) {
	return db.GetIDsWithContext(target, noCTX, ids...)
}

func (db DB) GetIDsWithContext(target interface{}, ctx map[string]interface{}, ids ...gouuidv6.UUID) (int, error) {
	t := time.Now()

	txn := db.KV.NewTransaction(false)
	defer txn.Discard()

	n, err := db.getIDsWithContext(txn, target, ctx, ids...)

	if db.Options.DebugMode {
		debugLogGet(target, t, n, err, ids...)
	}

	return n, err
}
