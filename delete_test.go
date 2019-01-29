package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Test_Delete(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	fullStruct := testtypes.FullStruct{}

	db.Save(&fullStruct)

	// Test the the fullStruct has been saved
	retrievedFullStruct := testtypes.FullStruct{}
	ok, _ := db.Get(&retrievedFullStruct, fullStruct.ID)
	if !ok || fullStruct.ID != retrievedFullStruct.ID {
		t.Error("Testing delete. Test fullStruct not saved correctly")
	}

	// Delete
	n, err := db.Delete("fullstruct", fullStruct.ID)

	if err != nil {
		t.Errorf("Testing delete. Got error %v", err)
	}

	if n != 1 {
		t.Errorf("Testing delete. Expected n = 1, got n = %v", n)
	}

	// Attempt to retrieve again
	ok, _ = db.Get(&retrievedFullStruct, fullStruct.ID)
	if ok {
		t.Error("Testing delete. Supposedly deleted fullStruct found on 2nd get")
	}
}

func Test_Delete_Multiple(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	fullStruct1 := testtypes.FullStruct{}
	fullStruct2 := testtypes.FullStruct{}
	fullStruct3 := testtypes.FullStruct{}

	db.Save(&fullStruct1, &fullStruct2, &fullStruct3)

	// Delete
	n, err := db.Delete("fullstruct", fullStruct1.ID, fullStruct2.ID, fullStruct3.ID)

	if err != nil {
		t.Errorf("Testing multiple delete. Got error %v", err)
	}

	if n != 3 {
		t.Errorf("Testing multiple delete. Expected n = %v, got n = %v", 3, n)
	}

	var fullStructs []testtypes.FullStruct
	c, _ := db.Find(&fullStructs).Count()
	if c > 0 {
		t.Errorf("Testing delete. Should have found any fullStructs, but found %v", c)
	}
}
