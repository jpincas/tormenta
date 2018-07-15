package tormentadb_test

import (
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// Index Creation
func Test_CreateIndexKeys(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// Create basic order and save
	// Orders have an 'index' on Customer field
	order := demo.Order{
		Department:  99,
		Customer:    "jon",
		ShippingFee: 5.99,
	}

	db.Save(&order)

	db.KV.View(func(txn *badger.Txn) error {
		customerIndex := tormenta.IndexKey([]byte("order"), order.ID, "customer", "jon")
		departmentIndex := tormenta.IndexKey([]byte("order"), order.ID, "department", 99)
		shippingFeeIndex := tormenta.IndexKey([]byte("order"), order.ID, "shippingfee", 5.99)

		_, err := txn.Get(customerIndex)
		if err == badger.ErrKeyNotFound {
			t.Error("Testing basic index key setting (string). Could not get index key")
		}

		_, err = txn.Get(departmentIndex)
		if err == badger.ErrKeyNotFound {
			t.Error("Testing basic index key setting (int). Could not get index key")
		}

		_, err = txn.Get(shippingFeeIndex)
		if err == badger.ErrKeyNotFound {
			t.Error("Testing basic index key setting (float64). Could not get index key")
		}

		return nil
	})
}
