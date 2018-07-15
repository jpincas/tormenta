package tormentadb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

var zeroValueTime time.Time

func Test_BasicSave(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// Create basic order and save
	order := demo.Order{}
	n, err := db.Save(&order)

	// Test any error
	if err != nil {
		t.Errorf("Testing basic record save. Got error %v", err)
	}

	// Test that 1 record was reported saved
	if n != 1 {
		t.Errorf("Testing basic record save. Expected 1 record saved, got %v", n)
	}

	// Check ID has been set
	if order.ID.IsNil() {
		t.Error("Testing basic record save with create new ID. ID after save is nil")
	}

	//  Check that updated field was set
	if order.LastUpdated == zeroValueTime {
		t.Error("Testing basic record save. 'Last Upated' is time zero value")
	}

	// Take a snapshot
	orderBeforeSecondSave := order

	// Save again
	n2, err2 := db.Save(&order)

	// Basic tests
	if err2 != nil {
		t.Errorf("Testing 2nd record save. Got error %v", err)
	}

	if n2 != 1 {
		t.Errorf("Testing 2nd record save. Expected 1 record saved, got %v", n)
	}

	//  Check that updated field was updated:the new value
	// should obviously be later
	if !orderBeforeSecondSave.LastUpdated.Before(order.LastUpdated) {
		t.Error("Testing 2nd record save. 'Created' time has changed")
	}
}

func Test_SaveTrigger(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// Create basic order and save
	order := demo.Order{}
	db.Save(&order)

	// Test postsave trigger
	if !order.OrderSaved {
		t.Error("Testing postsave trigger.  OrderSaved should be true but was not")
	}

	// Set up a condition that will cause the order not to save
	order.ContainsProhibitedItems = true

	// Test presave trigger
	n, err := db.Save(&order)
	if n != 0 || err == nil {
		t.Error("Testing presave trigger.  This record should not have saved, but it did and no error returned")
	}

}

func Test_SaveMultiple(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	order1 := demo.Order{}
	order2 := demo.Order{}

	// Multiple argument syntax
	n, _ := db.Save(&order1, &order2)
	if n != 2 {
		t.Errorf("Testing multiple save. Expected %v, got %v", 2, n)
	}

	// Spread syntax
	// A little akward as you can't just pass in the slice of entities
	// You have to manually translate to []Tormentable
	var ordersToSave []tormenta.Tormentable
	orders := []demo.Order{order1, order2}

	for _, order := range orders {
		ordersToSave = append(ordersToSave, &order)
	}

	n, _ = db.Save(ordersToSave...)
	if n != 2 {
		t.Errorf("Testing multiple save. Expected %v, got %v", 2, n)
	}

}

func Test_SaveMultipleLarge(t *testing.T) {
	const noOrders = 100000

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var ordersToSave []tormenta.Tormentable

	for i := 0; i < noOrders; i++ {
		ordersToSave = append(ordersToSave, &demo.Order{
			Customer: fmt.Sprintf("customer-%v", i),
		})
	}

	n, _ := db.Save(ordersToSave...)
	if n != noOrders {
		t.Errorf("Testing save large number of entities. Expected %v, got %v", noOrders, n)
	}

	var orders []demo.Order
	n, _ = db.Find(&orders).Run()
	if n != noOrders {
		t.Errorf("Testing save large number of entities, then retrieve. Expected %v, got %v", noOrders, n)
	}

}

func Test_SaveMultipleMillion(t *testing.T) {
	const noOrders = 1000000

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var ordersToSave []tormenta.Tormentable

	for i := 0; i < noOrders; i++ {
		ordersToSave = append(ordersToSave, &demo.Order{
			Customer: fmt.Sprintf("customer-%v", i),
		})
	}

	n, _ := db.Save(ordersToSave...)
	if n != noOrders {
		t.Errorf("Testing save large number of entities. Expected %v, got %v", noOrders, n)
	}

	var orders []demo.Order
	n, _ = db.Find(&orders).Run()
	if n != noOrders {
		t.Errorf("Testing save large number of entities, then retrieve. Expected %v, got %v", noOrders, n)
	}

}
