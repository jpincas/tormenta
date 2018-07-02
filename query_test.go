package tormenta

import (
	"testing"
)

func Test_BasicQuery(t *testing.T) {
	db, _ := OpenTest("data/tests")
	defer db.Close()

	// 1 order
	order1 := Order{}
	db.Save(&order1)

	var orders []Order
	n, err := db.Query(&orders).Run()

	if err != nil {
		t.Error("Testing basic querying - got error")
	}

	if len(orders) != 1 || n != 1 {
		t.Errorf("Testing querying with 1 entity saved. Expecting 1 entity - got %v/%v", len(orders), n)
	}

	// 2 orders
	order2 := Order{}
	orders = []Order{}
	db.Save(&order2)

	n, _ = db.Query(&orders).Run()

	if len(orders) != 2 || n != 2 {
		t.Errorf("Testing querying with 2 entity saved. Expecting 2 entities - got %v/%v", len(orders), n)
	}

}
