package tormentadb_test

import (
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
	"github.com/jpincas/tormenta/tormentadb/types"
)

// Index Creation
func Test_CreateIndexKeys(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// Create basic order and save
	// Orders have an 'index' on Customer field
	entity := types.TestType{
		IntField:                1,
		StringField:             "test",
		FloatField:              0.99,
		BoolField:               true,
		IntSliceField:           []int{1, 2},
		StringSliceField:        []string{"test1", "test2"},
		FloatSliceField:         []float64{0.99, 1.99},
		BoolSliceField:          []bool{true, false},
		DefinedIntField:         types.DefinedInt(1),
		DefinedStringField:      types.DefinedString("test"),
		DefinedFloatField:       types.DefinedFloat(0.99),
		DefinedBoolField:        types.DefinedBool(true),
		DefinedIntSliceField:    []types.DefinedInt{1, 2},
		DefinedStringSliceField: []types.DefinedString{"test1", "test2"},
		DefinedFloatSliceField:  []types.DefinedFloat{0.99, 1.99},
		DefinedBoolSliceField:   []types.DefinedBool{true, false},
	}

	db.Save(&entity)

	testCases := []struct {
		testName   string
		indexName  string
		indexValue interface{}
	}{
		// Basic types
		{"int field", "intfield", 1},
		{"string field", "stringfield", "test"},
		{"float field", "floatfield", 0.99},
		{"bool field", "boolfield", true},

		// Slice types - check both members
		{"int slice field", "intslicefield", 1},
		{"int slice field", "intslicefield", 2},
		{"string slice field", "stringslicefield", "test1"},
		{"string slice field", "stringslicefield", "test2"},
		{"float slice field", "floatslicefield", 0.99},
		{"float slice field", "floatslicefield", 1.99},
		{"bool slice field", "boolslicefield", true},
		{"bool slice field", "boolslicefield", false},

		// Defined types
		{"defined int field", "definedintfield", 1},
		{"defined string field", "definedstringfield", "test"},
		{"defined float field", "definedfloatfield", 0.99},
		{"defined bool field", "definedboolfield", true},
	}

	db.KV.View(func(txn *badger.Txn) error {

		for _, testCase := range testCases {
			i := tormenta.IndexKey([]byte("testtype"), entity.ID, testCase.indexName, testCase.indexValue)

			_, err := txn.Get(i)
			if err == badger.ErrKeyNotFound {
				t.Errorf("Testing %s. Could not get index key", testCase.testName)
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

	// content words
	expectedKeys := [][]byte{
		tormenta.IndexKey([]byte("product"), product.ID, "name", "coolest"),
		tormenta.IndexKey([]byte("product"), product.ID, "name", "product"),
		tormenta.IndexKey([]byte("product"), product.ID, "name", "world"),
	}

	// non content words
	nonExpectedKeys := [][]byte{
		tormenta.IndexKey([]byte("product"), product.ID, "name", "the"),
		tormenta.IndexKey([]byte("product"), product.ID, "name", "in"),
	}

	db.KV.View(func(txn *badger.Txn) error {
		for _, key := range expectedKeys {
			_, err := txn.Get(key)
			if err == badger.ErrKeyNotFound {
				t.Errorf("Testing index creation from slices.  Key [%v] should have been created but could not be retrieved", key)
			}
		}

		for _, key := range nonExpectedKeys {
			_, err := txn.Get(key)
			if err != badger.ErrKeyNotFound {
				t.Errorf("Testing index creation from slices.  Key [%v] should NOT have been created but was", key)
			}
		}

		return nil
	})
}
