// +build ignore

package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/demo"
)

func TestRandomise(t *testing.T) {
	// Make a list of 100 orders
	var orders []tormenta.Record
	for i := 0; i <= 100; i++ {
		orders = append(orders, &demo.Order{Department: i})
	}

	// Make a copy of the list before randomising, then randomise
	ordersBeforeRand := make([]tormenta.Record, len(orders))
	copy(ordersBeforeRand, orders)
	tormenta.RandomiseRecords(orders)

	// Go through element by element, compare, and set a flag to true if a difference was found
	foundDiff := false
	for i := range orders {
		if orders[i].(*demo.Order).Department != ordersBeforeRand[i].(*demo.Order).Department {
			foundDiff = true
		}
	}

	// If no differences were found, then fail
	if !foundDiff {
		t.Error("Testing randomise slice. Could not find any differences after randomisation")
	}

}
