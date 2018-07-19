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

func Test_CreateIndexKeys_Slice(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// Create basic order and save
	// Orders have an 'index' on Customer field
	product := demo.Product{
		Tags:        []string{"tag1", "tag2"},
		Departments: []int{1, 2, 3},
	}

	db.Save(&product)

	expectedKeys := [][]byte{
		tormenta.IndexKey([]byte("product"), product.ID, "tags", "tag1"),
		tormenta.IndexKey([]byte("product"), product.ID, "tags", "tag2"),
		tormenta.IndexKey([]byte("product"), product.ID, "departments", 1),
		tormenta.IndexKey([]byte("product"), product.ID, "departments", 2),
		tormenta.IndexKey([]byte("product"), product.ID, "departments", 3),
	}

	db.KV.View(func(txn *badger.Txn) error {
		for _, key := range expectedKeys {
			_, err := txn.Get(key)
			if err == badger.ErrKeyNotFound {
				t.Errorf("Testing index creation from slices.  Key [%v] should have been created but could not be retrieved", key)
			}
		}

		return nil
	})
}

func Test_CreateIndexKeys_Split(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	product := demo.Product{
		Name: "the coolest product in the world",
	}

	db.Save(&product)

	expectedKeys := [][]byte{
		tormenta.IndexKey([]byte("product"), product.ID, "name", "the"),
		tormenta.IndexKey([]byte("product"), product.ID, "name", "coolest"),
		tormenta.IndexKey([]byte("product"), product.ID, "name", "product"),
		tormenta.IndexKey([]byte("product"), product.ID, "name", "in"),
		tormenta.IndexKey([]byte("product"), product.ID, "name", "the"),
		tormenta.IndexKey([]byte("product"), product.ID, "name", "world"),
	}

	db.KV.View(func(txn *badger.Txn) error {
		for _, key := range expectedKeys {
			_, err := txn.Get(key)
			if err == badger.ErrKeyNotFound {
				t.Errorf("Testing index creation from slices.  Key [%v] should have been created but could not be retrieved", key)
			}
		}

		return nil
	})
}
