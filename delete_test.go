package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Test_Delete_EntityID(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	fullStruct := testtypes.FullStruct{}

	db.Save(&fullStruct)

	// Test the the fullStruct has been saved
	retrievedFullStruct := testtypes.FullStruct{}
	ok, _ := db.Get(&retrievedFullStruct, fullStruct.ID)
	if !ok || fullStruct.ID != retrievedFullStruct.ID {
		t.Error("Testing delete. Test fullStruct not saved correctly")
	}

	// Delete by entity id
	err := db.Delete(&fullStruct)

	if err != nil {
		t.Errorf("Testing delete. Got error %v", err)
	}

	// Attempt to retrieve again
	ok, _ = db.Get(&retrievedFullStruct, fullStruct.ID)
	if ok {
		t.Error("Testing delete. Supposedly deleted fullStruct found on 2nd get")
	}
}

func Test_Delete_SeparateID(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	fullStruct := testtypes.FullStruct{}
	fullStruct2 := testtypes.FullStruct{}

	db.Save(&fullStruct, &fullStruct2)

	// Test the the fullStruct has been saved
	retrievedFullStruct := testtypes.FullStruct{}
	ok, _ := db.Get(&retrievedFullStruct, fullStruct.ID)
	if !ok || fullStruct.ID != retrievedFullStruct.ID {
		t.Error("Testing delete. Test fullStruct not saved correctly")
	}

	// Test the the fullStruct has been saved
	retrievedFullStruct2 := testtypes.FullStruct{}
	ok, _ = db.Get(&retrievedFullStruct2, fullStruct2.ID)
	if !ok || fullStruct2.ID != retrievedFullStruct2.ID {
		t.Error("Testing delete. Test fullStruct not saved correctly")
	}

	// Delete by separate id
	// We're being tricky here - we're passing in the entity #2,
	// but specifying the ID of #1 to delete
	err := db.Delete(&fullStruct2, fullStruct.ID)

	if err != nil {
		t.Errorf("Testing delete. Got error %v", err)
	}

	// Attempt to retrieve again
	ok, _ = db.Get(&retrievedFullStruct, fullStruct.ID)
	if ok {
		t.Error("Testing delete. Supposedly deleted fullStruct found on 2nd get")
	}
}
