package tormenta_test

import (
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Test_BasicGet(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	// Create basic fullStruct and save, then blank the ID
	fullStruct := testtypes.FullStruct{}

	if _, err := db.Save(&fullStruct); err != nil {
		t.Errorf("Testing get entity without ID. Got error on save (%v)", err)
	}

	ttIDBeforeBlanking := fullStruct.ID
	fullStruct.ID = gouuidv6.UUID{}

	// Attempt to get entity without ID
	found, err := db.Get(&fullStruct)
	if err != nil {
		t.Errorf("Testing get entity without ID. Got error (%v) but should simply fail to find", err)
	}

	if found {
		t.Errorf("Testing get entity without ID. Expected not to find anything, but did")

	}

	// Reset the fullStruct ID
	fullStruct.ID = ttIDBeforeBlanking
	ok, err := db.Get(&fullStruct)
	if err != nil {
		t.Errorf("Testing basic record get. Got error %v", err)
	}

	if !ok {
		t.Error("Testing basic record get. Record was not found")
	}

}

func Test_GetByID(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	fullStruct := testtypes.FullStruct{}
	tt2 := testtypes.FullStruct{}
	db.Save(&fullStruct)

	// Overwite ID
	ok, err := db.Get(&tt2, fullStruct.ID)

	if err != nil {
		t.Errorf("Testing get by id. Got error %v", err)
	}

	if !ok {
		t.Error("Testing get by id. Record was not found")
	}

	if fullStruct.ID != tt2.ID {
		t.Error("Testing get by id. Entity retreived by ID was not the same as that saved")
	}
}

func Test_GetByMultipleIDs(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	noOfTests := 500

	var toSave []tormenta.Record
	var ids []gouuidv6.UUID

	for i := 0; i < noOfTests; i++ {
		id := gouuidv6.NewFromTime(time.Now())
		record := testtypes.FullStruct{}
		record.SetID(id)
		toSave = append(toSave, &record)
		ids = append(ids, id)
	}

	if _, err := db.Save(toSave...); err != nil {
		t.Errorf("Testing get by multiple ids. Got error saving %v", err)
	}

	var results []testtypes.FullStruct
	n, err := db.GetIDs(&results, ids...)

	if err != nil {
		t.Errorf("Testing get by multiple ids. Got error %v", err)
	}

	if n != len(results) {
		t.Errorf("Testing get by multiple ids. Mismatch between reported n (%v) and length of results slice (%v)", n, len(results))
	}

	if n != len(ids) {
		t.Errorf("Testing get by multiple ids. Wanted %v results, got %v", len(ids), n)
	}

	for i, _ := range results {
		if results[i].ID != toSave[i].GetID() {
			t.Errorf("Testing get by multiple ids. ID mismatch for array member %v. Wanted %v, got %v", i, toSave[i].GetID(), results[i].ID)
		}
	}
}

func Test_GetTriggers(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	fullStruct := testtypes.FullStruct{}
	db.Save(&fullStruct)
	ok, err := db.Get(&fullStruct)

	if err != nil {
		t.Errorf("Testing get triggers. Got error %v", err)
	}

	if !ok {
		t.Error("Testing get triggers. Record was not found")
	}

	if !fullStruct.Retrieved {
		t.Error("Testing get triggers.  Expected ttRetrieved = true; got false")
	}
}
