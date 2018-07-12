package tormenta

import (
	"fmt"
	"testing"

	"github.com/jpincas/gouuidv6"
)

func Test_BasicGet(t *testing.T) {
	db, _ := OpenTest("data/tests")
	defer db.Close()

	// Test struc that has no model
	noModel := NoModel{}
	keyRoot, _ := entityTypeAndValue(&noModel)
	_, err := db.Get(&noModel)

	if err == nil || err.Error() != fmt.Sprintf(errNoModel, keyRoot) {
		t.Error("Testing save entity without model. Produced wrong error or no error")
	}

	// Create basic order and save, then blank the ID
	order := Order{}
	keyRoot, _ = entityTypeAndValue(&order)
	db.Save(&order)
	orderIDBeforeBlanking := order.ID
	order.ID = gouuidv6.UUID{}

	// Attempt to get entity without ID
	_, err = db.Get(&order)
	if err == nil || err.Error() != fmt.Sprintf(errNoID, keyRoot) {
		t.Error("Testing get entity without ID. Produced wrong error or no error")
	}

	// Reset the order ID
	order.ID = orderIDBeforeBlanking
	ok, err := db.Get(&order)
	if err != nil {
		t.Errorf("Testing basic record get. Got error %v", err)
	}

	if !ok {
		t.Error("Testing basic record get. Record was not found")
	}

}

func Test_GetByID(t *testing.T) {
	db, _ := OpenTest("data/tests")
	defer db.Close()

	order := Order{}
	order2 := Order{}
	db.Save(&order)

	// Overwite ID
	ok, err := db.Get(&order2, order.ID)

	if err != nil {
		t.Errorf("Testing get by id. Got error %v", err)
	}

	if !ok {
		t.Error("Testing get by id. Record was not found")
	}

	if order.ID != order2.ID {
		t.Error("Testing get by id. Entity retreived by ID was not the same as that saved")
	}

}
