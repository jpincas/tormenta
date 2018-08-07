package tormentadb_test

import (
	"testing"
	"time"

	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

func Test_BasicQuery(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// 1 order
	order1 := demo.Order{}
	db.Save(&order1)

	var orders []demo.Order
	n, _, err := db.Find(&orders).Run()

	if err != nil {
		t.Error("Testing basic querying - got error")
	}

	if len(orders) != 1 || n != 1 {
		t.Errorf("Testing querying with 1 entity saved. Expecting 1 entity - got %v/%v", len(orders), n)
	}

	orders = []demo.Order{}
	c, _, err := db.Find(&orders).Count()
	if c != 1 {
		t.Errorf("Testing count 1 entity saved. Expecting 1 - got %v", c)
	}

	// 2 orders
	order2 := demo.Order{}
	db.Save(&order2)

	orders = []demo.Order{}
	if n, _, _ := db.Find(&orders).Run(); n != 2 {
		t.Errorf("Testing querying with 2 entity saved. Expecting 2 entities - got %v", n)
	}

	if c, _, _ := db.Find(&orders).Count(); c != 2 {
		t.Errorf("Testing count 2 entities saved. Expecting 2 - got %v", c)
	}
	if order1.ID == order2.ID {
		t.Errorf("Testing querying with 2 entities saved. 2 entities saved both have same ID")
	}
	if orders[0].ID == orders[1].ID {
		t.Errorf("Testing querying with 2 entities saved. 2 results returned. Both have same ID")
	}

	// Limit
	orders = []demo.Order{}
	if n, _, _ := db.Find(&orders).Limit(1).Run(); n != 1 {
		t.Errorf("Testing querying with 2 entities saved + limit. Wrong number of results received")
	}

	// Reverse - simple, only tests number received
	orders = []demo.Order{}
	if n, _, _ := db.Find(&orders).Reverse().Run(); n != 2 {
		t.Errorf("Testing querying with 2 entities saved + reverse. Expected %v, got %v", 2, n)
	}

	// Reverse + Limit - simple, only tests number received
	orders = []demo.Order{}
	if n, _, _ := db.Find(&orders).Reverse().Limit(1).Run(); n != 1 {
		t.Errorf("Testing querying with 2 entities saved + reverse + limit. Expected %v, got %v", 1, n)
	}

	// Reverse + Count
	orders = []demo.Order{}
	if c, _, _ := db.Find(&orders).Reverse().Count(); c != 2 {
		t.Errorf("Testing count with 2 entities saved + reverse. Expected %v, got %v", 2, c)
	}

}

func Test_BasicQuery_First(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	order1 := demo.Order{}
	order2 := demo.Order{}
	db.Save(&order1, &order2)

	var order demo.Order
	n, _, err := db.First(&order).Run()

	if err != nil {
		t.Error("Testing first - got error")
	}

	if n != 1 {
		t.Errorf("Testing first. Expecting 1 entity - got %v", n)
	}

	if order.ID.IsNil() {
		t.Errorf("Testing first. Nil ID retrieved")
	}

	if order.ID != order1.ID {
		t.Errorf("Testing first. Order IDs are not equal - wrong order retrieved")
	}

	// Test nothing found (impossible range)
	n, _, _ = db.First(&order).From(time.Now()).To(time.Now()).Run()
	if n != 0 {
		t.Errorf("Testing first when nothing should be found.  Got n = %v", n)
	}
}
