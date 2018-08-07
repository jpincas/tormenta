package tormentadb_test

import (
	"fmt"
	"testing"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

func Test_BasicGet(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// Create basic order and save, then blank the ID
	order := demo.Order{}
	keyRoot := tormenta.KeyRoot(&order)
	db.Save(&order)
	orderIDBeforeBlanking := order.ID
	order.ID = gouuidv6.UUID{}

	// Attempt to get entity without ID
	_, _, err := db.Get(&order)
	if err == nil || err.Error() != fmt.Sprintf(tormenta.ErrNoID, keyRoot) {
		t.Error("Testing get entity without ID. Produced wrong error or no error")
	}

	// Reset the order ID
	order.ID = orderIDBeforeBlanking
	ok, _, err := db.Get(&order)
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

	order := demo.Order{}
	order2 := demo.Order{}
	db.Save(&order)

	// Overwite ID
	ok, _, err := db.Get(&order2, order.ID)

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

func Test_GetTriggers(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	order := demo.Order{}
	db.Save(&order)
	ok, _, err := db.Get(&order)

	if err != nil {
		t.Errorf("Testing get triggers. Got error %v", err)
	}

	if !ok {
		t.Error("Testing get triggers. Record was not found")
	}

	if !order.OrderRetrieved {
		t.Error("Testing get triggers.  Expected OrderRetrieved = true; got false")
	}
}
