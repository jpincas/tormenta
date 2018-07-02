package tormenta

import (
	"testing"
)

func Test_BasicQuery(t *testing.T) {
	db, _ := OpenTest("data/tests")
	defer db.Close()

	// Save a couple of orders
	order1 := Order{}
	order2 := Order{}
	db.Save(&order1, &order2)

	// Suboptimal:

	// Run a query for all orders
	var order Order
	results, _ := db.Query(&order).Run()

	// Conver the []tormentable -> []orders
	var orders []Order
	for _, result := range results {
		orders = append(orders, *result.(*Order))
	}

	// What I really want to do:
	// var orders []Order
	// db.Query(&orders).Run()

}
