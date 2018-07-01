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
	keyRoot, _ := getKeyRoot(&noModel)
	_, err := db.Get(&noModel)

	if err == nil || err.Error() != fmt.Sprintf(errNoModel, keyRoot) {
		t.Error("Testing save entity without model. Produced wrong error or no error")
	}

	// Create basic order and save, then blank the ID
	order := Order{}
	keyRoot, _ = getKeyRoot(&order)
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
