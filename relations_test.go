package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
)

func Test_Relations_Load(t *testing.T) {
	// Open the DB
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	// Create some products
	product1 := Product{
		Code:          "SKU1",
		Name:          "Product1",
		Price:         1.00,
		StartingStock: 50}
	product2 := Product{
		Code:          "SKU2",
		Name:          "Product2",
		Price:         2.00,
		StartingStock: 100}

	// Save them
	db.Save(&product1, &product2)

	// Create some orders for those products
	// 1 product per order for now
	order1 := Order{
		Customer:    "Mr T",
		Department:  1,
		ShippingFee: 4.99,
		ProductID:   product1.ID,
	}

	order2 := Order{
		Customer:    "Mr T",
		Department:  1,
		ShippingFee: 4.99,
		ProductID:   product2.ID,
	}

	order3 := Order{
		Customer:    "Mr T",
		Department:  1,
		ShippingFee: 4.99,
		ProductID:   product1.ID,
	}

	// Save
	db.Save(&order1, &order2, &order3)

	// Reload
	var orders []Order
	if n, err := db.Find(&orders).Run(); err != nil || n != 3 {
		t.Errorf("Save/retrieve failed. Err: %v; n: %v", err, n)
	}

	// Convert to records
	var records []tormenta.Record
	for i := range orders {
		records = append(records, &orders[i])
	}

	// Attempt to load relations
	if err := tormenta.LoadRelations(db, "Product", records...); err != nil {
		t.Errorf("Error loading relations: %s", err)
	}

	for i, order := range orders {
		if order.ProductID != order.Product.ID {
			t.Errorf(
				"Comparing ProductID to Product.ID for order %v and they are not the same: %v vs %v",
				i,
				order.ProductID,
				order.Product.ID,
			)
		}
	}
}
