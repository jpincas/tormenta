package tormenta_test

import (
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

// Index Creation
func Test_MakeIndexKeys(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	entity := testtypes.FullStruct{
		IntField:                1,
		StringField:             "test",
		FloatField:              0.99,
		BoolField:               true,
		IntSliceField:           []int{1, 2},
		StringSliceField:        []string{"test1", "test2"},
		FloatSliceField:         []float64{0.99, 1.99},
		BoolSliceField:          []bool{true, false},
		DefinedIntField:         testtypes.DefinedInt(1),
		DefinedStringField:      testtypes.DefinedString("test"),
		DefinedFloatField:       testtypes.DefinedFloat(0.99),
		DefinedBoolField:        testtypes.DefinedBool(true),
		DefinedIntSliceField:    []testtypes.DefinedInt{1, 2},
		DefinedStringSliceField: []testtypes.DefinedString{"test1", "test2"},
		DefinedFloatSliceField:  []testtypes.DefinedFloat{0.99, 1.99},
		DefinedBoolSliceField:   []testtypes.DefinedBool{true, false},
		MyStruct: testtypes.MyStruct{
			StructIntField:    1,
			StructStringField: "test",
			StructFloatField:  0.99,
			StructBoolField:   true,
		},
	}

	db.Save(&entity)

	testCases := []struct {
		testName   string
		indexName  string
		indexValue interface{}
	}{
		// Basic testtypes
		{"int field", "intfield", 1},
		{"string field", "stringfield", "test"},
		{"float field", "floatfield", 0.99},
		{"bool field", "boolfield", true},

		// Slice testtypes - check both members
		{"int slice field", "intslicefield", 1},
		{"int slice field", "intslicefield", 2},
		{"string slice field", "stringslicefield", "test1"},
		{"string slice field", "stringslicefield", "test2"},
		{"float slice field", "floatslicefield", 0.99},
		{"float slice field", "floatslicefield", 1.99},
		{"bool slice field", "boolslicefield", true},
		{"bool slice field", "boolslicefield", false},

		// Defined testtypes
		{"defined int field", "definedintfield", 1},
		{"defined string field", "definedstringfield", "test"},
		{"defined float field", "definedfloatfield", 0.99},
		{"defined bool field", "definedboolfield", true},

		// Struct structs
		{"embedded struct - int field", "structintfield", 1},
		{"embedded struct - string field", "structstringfield", "test"},
		{"embedded struct - float field", "structfloatfield", 0.99},
		{"embedded struct - bool field", "structboolfield", true},
	}

	// Step 1 - make sure that the keys that we expect are present after saving
	db.KV.View(func(txn *badger.Txn) error {

		for _, testCase := range testCases {
			i := tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, testCase.indexName, testCase.indexValue)

			_, err := txn.Get(i)
			if err == badger.ErrKeyNotFound {
				t.Errorf("Testing %s. Could not get index key", testCase.testName)
			}
		}

		return nil
	})

	// Step 2 - delete the record and test that it has been deindexed
	err := db.Delete(&entity)

	if err != nil {
		t.Errorf("Testing delete. Got error %v", err)
	}

	db.KV.View(func(txn *badger.Txn) error {

		for _, testCase := range testCases {
			i := tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, testCase.indexName, testCase.indexValue)

			if _, err := txn.Get(i); err != badger.ErrKeyNotFound {
				t.Errorf("Testing %s after deletion. Should not find index key but did", testCase.testName)
			}
		}

		return nil
	})
}

func Test_ReIndex(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	entity := testtypes.FullStruct{
		IntField:    1,
		StringField: "test",
	}

	// Save the entity first
	db.Save(&entity)

	// Step 1 - test that the 2 basic indexes have been created
	db.KV.View(func(txn *badger.Txn) error {
		key := tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, "intField", 1)
		if _, err := txn.Get(key); err == badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Could not get index key", "int field indexing")
		}

		key = tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, "stringField", "test")
		if _, err := txn.Get(key); err == badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Could not get index key", "string field indexing")
		}

		return nil
	})

	// Stpe 2 - Now make some changes and update
	entity.IntField = 2
	entity.StringField = "test_update"
	db.Save(&entity)

	// Step 3 - test that the 2 previous indices are gone
	db.KV.View(func(txn *badger.Txn) error {
		key := tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, "intField", 1)
		if _, err := txn.Get(key); err != badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Found index key when shouldn't have", "int field indexing")
		}

		key = tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, "stringField", "test")
		if _, err := txn.Get(key); err != badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Found index key when shouldn't have", "string field indexing")
		}

		return nil
	})

	// Step 4 - test that the 2 new indices are present
	db.KV.View(func(txn *badger.Txn) error {
		key := tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, "intField", 2)
		if _, err := txn.Get(key); err == badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Could not get index key after update", "int field indexing")
		}

		key = tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, "stringField", "test_update")
		if _, err := txn.Get(key); err == badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Could not get index key after update", "string field indexing")
		}

		return nil
	})

}

func Test_MakeIndexKeys_Split(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	fullStruct := testtypes.FullStruct{
		MultipleWordField: "the coolest fullStruct in the world",
	}

	db.Save(&fullStruct)

	// content words
	expectedKeys := [][]byte{
		tormenta.MakeIndexKey([]byte("fullstruct"), fullStruct.ID, "multiplewordfield", "coolest"),
		tormenta.MakeIndexKey([]byte("fullstruct"), fullStruct.ID, "multiplewordfield", "fullStruct"),
		tormenta.MakeIndexKey([]byte("fullstruct"), fullStruct.ID, "multiplewordfield", "world"),
	}

	// non content words
	nonExpectedKeys := [][]byte{
		tormenta.MakeIndexKey([]byte("fullstruct"), fullStruct.ID, "multiplewordfield", "the"),
		tormenta.MakeIndexKey([]byte("fullstruct"), fullStruct.ID, "multiplewordfield", "in"),
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
