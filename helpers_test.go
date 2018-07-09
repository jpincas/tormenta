package tormenta

import "testing"

func TestRandomise(t *testing.T) {
	// Make a list of 100 orders
	var orders []Tormentable
	for i := 0; i <= 100; i++ {
		orders = append(orders, &Order{Department: i})
	}

	// Make a copy of the list before randomising, then randomise
	ordersBeforeRand := make([]Tormentable, len(orders))
	copy(ordersBeforeRand, orders)
	randomiseTormentables(orders)

	// Go through element by element, compare, and set a flag to true if a difference was found
	foundDiff := false
	for i := range orders {
		if orders[i].(*Order).Department != ordersBeforeRand[i].(*Order).Department {
			foundDiff = true
		}
	}

	// If no differences were found, then fail
	if !foundDiff {
		t.Error("Testing randomise slice. Could not find any differences after randomisation")
	}

}
