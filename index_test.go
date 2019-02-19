package tormenta_test

import (
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

// Index Creation
func Test_MakeIndexKeys(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	id := gouuidv6.New()

	entity := testtypes.FullStruct{
		IntField:                1,
		IDField:                 id,
		StringField:             "test",
		FloatField:              0.99,
		BoolField:               true,
		DateField:               time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
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
		NoSaveSimple: "dontsaveitsodontindexit",
		StructField: testtypes.MyStruct{
			StructStringField: "test",
		},
	}

	db.Save(&entity)

	testCases := []struct {
		testName    string
		indexName   string
		indexValue  interface{}
		shouldIndex bool
	}{
		// Basic testtypes
		{"int field", "intfield", 1, true},
		{"id field", "idfield", id, true},
		{"string field", "stringfield", "test", true},
		{"float field", "floatfield", 0.99, true},
		{"bool field", "boolfield", true, true},
		{"date field", "datefield", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC).Unix(), true},

		// Slice testtypes - check both members
		{"int slice field", "intslicefield", 1, true},
		{"int slice field", "intslicefield", 2, true},
		{"string slice field", "stringslicefield", "test1", true},
		{"string slice field", "stringslicefield", "test2", true},
		{"float slice field", "floatslicefield", 0.99, true},
		{"float slice field", "floatslicefield", 1.99, true},
		{"bool slice field", "boolslicefield", true, true},
		{"bool slice field", "boolslicefield", false, true},

		// Defined testtypes
		{"defined int field", "definedintfield", 1, true},
		{"defined string field", "definedstringfield", "test", true},
		{"defined float field", "definedfloatfield", 0.99, true},
		{"defined bool field", "definedboolfield", true, true},

		// Anonymous structs
		{"embedded struct - int field", "structintfield", 1, true},
		{"embedded struct - string field", "structstringfield", "test", true},
		{"embedded struct - float field", "structfloatfield", 0.99, true},
		{"embedded struct - bool field", "structboolfield", true, true},

		// Names structs
		{"named struct - string field", "structfield.structstringfield", "test", true},

		// No save / No Index
		{"no index field simple", "noindexsimple", "dontsaveitsodontindexit", false},
		{"no index field two tags", "noindextwotags", "dontsaveitsodontindexit", false},
		{"no index field, two tags, different order", "noindextwotagsdifferentorder", "dontsaveitsodontindexit", false},
		{"no save field", "nosavesimple", "dontsaveitsodontindexit", false},
	}

	// Step 1 - make sure that the keys that we expect are present after saving
	db.KV.View(func(txn *badger.Txn) error {

		for _, testCase := range testCases {
			i := tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, []byte(testCase.indexName), testCase.indexValue)

			_, err := txn.Get(i)
			if testCase.shouldIndex && err == badger.ErrKeyNotFound {
				t.Errorf("Testing %s. Could not get index key", testCase.testName)
			} else if !testCase.shouldIndex && err != badger.ErrKeyNotFound {
				t.Errorf("Testing %s. Should not have found the index key but did", testCase.testName)
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
			i := tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, []byte(testCase.indexName), testCase.indexValue)

			if _, err := txn.Get(i); err != badger.ErrKeyNotFound {
				t.Errorf("Testing %s after deletion. Should not find index key but did", testCase.testName)
			}
		}

		return nil
	})
}

func Test_ReIndex(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	entity := testtypes.FullStruct{
		IntField:    1,
		StringField: "test",
	}

	// Save the entity first
	db.Save(&entity)

	// Step 1 - test that the 2 basic indexes have been created
	db.KV.View(func(txn *badger.Txn) error {
		key := tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, []byte("intField"), 1)
		if _, err := txn.Get(key); err == badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Could not get index key", "int field indexing")
		}

		key = tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, []byte("stringField"), "test")
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
		key := tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, []byte("intField"), 1)
		if _, err := txn.Get(key); err != badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Found index key when shouldn't have", "int field indexing")
		}

		key = tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, []byte("stringField"), "test")
		if _, err := txn.Get(key); err != badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Found index key when shouldn't have", "string field indexing")
		}

		return nil
	})

	// Step 4 - test that the 2 new indices are present
	db.KV.View(func(txn *badger.Txn) error {
		key := tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, []byte("intField"), 2)
		if _, err := txn.Get(key); err == badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Could not get index key after update", "int field indexing")
		}

		key = tormenta.MakeIndexKey([]byte("fullstruct"), entity.ID, []byte("stringField"), "test_update")
		if _, err := txn.Get(key); err == badger.ErrKeyNotFound {
			t.Errorf("Testing %s. Could not get index key after update", "string field indexing")
		}

		return nil
	})

}

func Test_MakeIndexKeys_Split(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	fullStruct := testtypes.FullStruct{
		MultipleWordField: "the coolest fullStruct in the world",
	}

	db.Save(&fullStruct)

	// content words
	expectedKeys := [][]byte{
		tormenta.MakeIndexKey([]byte("fullstruct"), fullStruct.ID, []byte("multiplewordfield"), "coolest"),
		tormenta.MakeIndexKey([]byte("fullstruct"), fullStruct.ID, []byte("multiplewordfield"), "fullStruct"),
		tormenta.MakeIndexKey([]byte("fullstruct"), fullStruct.ID, []byte("multiplewordfield"), "world"),
	}

	// non content words
	nonExpectedKeys := [][]byte{
		tormenta.MakeIndexKey([]byte("fullstruct"), fullStruct.ID, []byte("multiplewordfield"), "the"),
		tormenta.MakeIndexKey([]byte("fullstruct"), fullStruct.ID, []byte("multiplewordfield"), "in"),
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
