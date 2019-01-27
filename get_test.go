package tormenta_test

import (
	"fmt"
	"testing"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
)

func Test_BasicGet(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// Create basic tt and save, then blank the ID
	tt := TestType{}
	keyRoot := tormenta.KeyRoot(&tt)
	db.Save(&tt)
	ttIDBeforeBlanking := tt.ID
	tt.ID = gouuidv6.UUID{}

	// Attempt to get entity without ID
	_, _, err := db.Get(&tt)
	if err == nil || err.Error() != fmt.Sprintf(tormenta.ErrNoID, keyRoot) {
		t.Error("Testing get entity without ID. Produced wrong error or no error")
	}

	// Reset the tt ID
	tt.ID = ttIDBeforeBlanking
	ok, _, err := db.Get(&tt)
	if err != nil {
		t.Errorf("Testing basic record get. Got error %v", err)
	}

	if !ok {
		t.Error("Testing basic record get. Record was not found")
	}

}

func Test_GetByID(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	tt := TestType{}
	tt2 := TestType{}
	db.Save(&tt)

	// Overwite ID
	ok, _, err := db.Get(&tt2, tt.ID)

	if err != nil {
		t.Errorf("Testing get by id. Got error %v", err)
	}

	if !ok {
		t.Error("Testing get by id. Record was not found")
	}

	if tt.ID != tt2.ID {
		t.Error("Testing get by id. Entity retreived by ID was not the same as that saved")
	}
}

func Test_GetTriggers(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	tt := TestType{}
	db.Save(&tt)
	ok, _, err := db.Get(&tt)

	if err != nil {
		t.Errorf("Testing get triggers. Got error %v", err)
	}

	if !ok {
		t.Error("Testing get triggers. Record was not found")
	}

	if !tt.Retrieved {
		t.Error("Testing get triggers.  Expected ttRetrieved = true; got false")
	}
}
