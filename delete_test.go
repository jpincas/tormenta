// +build ignore

package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/demo"
)

func Test_Delete(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	order := demo.Order{}

	db.Save(&order)

	// Test the the order has been saved
	retrievedOrder := demo.Order{}
	ok, _, _ := db.Get(&retrievedOrder, order.ID)
	if !ok || order.ID != retrievedOrder.ID {
		t.Error("Testing delete. Test order not saved correctly")
	}

	// Delete
	n, err := db.Delete("order", order.ID)

	if err != nil {
		t.Errorf("Testing delete. Got error %v", err)
	}

	if n != 1 {
		t.Errorf("Testing delete. Expected n = 1, got n = %v", n)
	}

	// Attempt to retrieve again
	ok, _, _ = db.Get(&retrievedOrder, order.ID)
	if ok {
		t.Error("Testing delete. Supposedly deleted order found on 2nd get")
	}
}

func Test_Delete_Multiple(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	order1 := demo.Order{}
	order2 := demo.Order{}
	order3 := demo.Order{}

	db.Save(&order1, &order2, &order3)

	// Delete
	n, err := db.Delete("order", order1.ID, order2.ID, order3.ID)

	if err != nil {
		t.Errorf("Testing multiple delete. Got error %v", err)
	}

	if n != 3 {
		t.Errorf("Testing multiple delete. Expected n = %v, got n = %v", 3, n)
	}

	var orders []demo.Order
	c, _, _ := db.Find(&orders).Count()
	if c > 0 {
		t.Errorf("Testing delete. Should have found any orders, but found %v", c)
	}
}
