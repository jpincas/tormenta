package tormenta

import (
	"fmt"
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
)

// Basic Queries

func Test_BasicQuery(t *testing.T) {
	db, _ := OpenTest("data/tests")
	defer db.Close()

	// 1 order
	order1 := Order{}
	db.Save(&order1)

	var orders []Order
	n, err := db.Find(&orders).Run()

	if err != nil {
		t.Error("Testing basic querying - got error")
	}

	if len(orders) != 1 || n != 1 {
		t.Errorf("Testing querying with 1 entity saved. Expecting 1 entity - got %v/%v", len(orders), n)
	}

	c, err := db.Find(&orders).Count()
	if c != len(orders) {
		t.Errorf("Testing count 1 entity saved. Expecting 1 - got %v", c)
	}

	// 2 orders
	order2 := Order{}
	orders = []Order{}
	db.Save(&order2)

	n, _ = db.Find(&orders).Run()

	if len(orders) != 2 || n != 2 {
		t.Errorf("Testing querying with 2 entity saved. Expecting 2 entities - got %v/%v", len(orders), n)
	}

	c, err = db.Find(&orders).Count()
	if c != len(orders) {
		t.Errorf("Testing count 2 entities saved. Expecting 2 - got %v", c)
	}

}

func Test_BasicQuery_First(t *testing.T) {
	db, _ := OpenTest("data/tests")
	defer db.Close()

	order1 := Order{}
	order2 := Order{}
	db.Save(&order1, &order2)

	var order Order
	n, err := db.First(&order).Run()

	if err != nil {
		t.Error("Testing first - got error")
	}

	if n != 1 {
		t.Errorf("Testing first. Expecting 1 entity - got %v", n)
	}

	if order.ID != order1.ID {
		t.Errorf("Testing first. Order IDs are not equal - wrong order retrieved")
	}

	// Test nothing found (impossible range)
	n, _ = db.First(&order).From(time.Now()).To(time.Now()).Run()
	if n != 0 {
		t.Errorf("Testing first when nothing should be found.  Got n = %v", n)
	}
}

func Test_BasicQuery_DateRange(t *testing.T) {
	// Create a list of orders over a date range
	var orders []Tormentable
	dates := []time.Time{
		// Now
		time.Now(),

		// Over the last week
		time.Now().Add(-1 * 24 * time.Hour),
		time.Now().Add(-2 * 24 * time.Hour),
		time.Now().Add(-3 * 24 * time.Hour),
		time.Now().Add(-4 * 24 * time.Hour),
		time.Now().Add(-5 * 24 * time.Hour),
		time.Now().Add(-6 * 24 * time.Hour),
		time.Now().Add(-7 * 24 * time.Hour),

		// Specific years
		time.Date(2009, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2010, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2011, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2012, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2013, time.January, 1, 1, 0, 0, 0, time.UTC),
	}

	for _, date := range dates {
		orders = append(orders, &Order{
			Model: Model{
				ID: gouuidv6.NewFromTime(date),
			},
		})
	}

	// Save the orders
	db, _ := OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	// Also another entity, to make sure there is no crosstalk
	db.Save(&Product{
		Code:          "001",
		Name:          "Computer",
		Price:         999.99,
		StartingStock: 50,
		Description:   defaultDescription})

	// Quick check that all orders have saved correctly
	var results []Order
	n, _ := db.Find(&results).Run()

	if len(results) != len(orders) || n != len(orders) {
		t.Errorf("Testing range query. Haven't even got to ranges yet. Just basic query expected %v - got %v/%v", len(orders), len(results), n)
		t.FailNow()
	}

	// Range test cases
	testCases := []struct {
		testName  string
		from, to  time.Time
		expected  int
		includeTo bool
		limit     int
		reverse   bool
	}{
		{"from right now - no orders expected, no 'to'", time.Now(), time.Time{}, 0, false, 0, false},
		{"from beginning of time - all orders should be included, no 'to'", time.Time{}, time.Time{}, len(orders), false, 0, false},
		{"from 2014, no 'to'", time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), time.Time{}, 8, false, 0, false},
		{"from 1 hour ago, no 'to'", time.Now().Add(-1 * time.Hour), time.Time{}, 1, false, 0, false},
		{"from beginning of time to now - expect all", time.Time{}, time.Now(), len(orders), true, 0, false},
		{"from beginning of time to 2014 - expect 5", time.Time{}, time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), 5, true, 0, false},
		{"from beginning of time to an hour ago - expect all but 1", time.Time{}, time.Now().Add(-1 * time.Hour), len(orders) - 1, true, 0, false},
		{"from beginning of time - limit 1", time.Time{}, time.Time{}, 1, false, 1, false},
		{"from beginning of time - limit 10", time.Time{}, time.Time{}, 10, false, 10, false},

		{"reversed - from beginning of time", time.Time{}, time.Time{}, 0, false, 0, true},
		{"reverse - from now - no to", time.Now(), time.Time{}, len(orders), false, 0, true},
		{"reverse - from now to 2014 - expect 8", time.Now(), time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), 8, true, 0, true},
		{"reverse - from now to 2014 - limit 5 - expect 5", time.Now(), time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), 5, true, 5, true},
	}

	for _, testCase := range testCases {
		rangequeryResults := []Order{}
		query := db.Find(&rangequeryResults).From(testCase.from)

		if testCase.includeTo {
			query = query.To(testCase.to)
		}

		if testCase.limit > 0 {
			query = query.Limit(testCase.limit)
		}

		if testCase.reverse {
			query = query.Reverse()
		}

		n, _ := query.Run()
		c, _ := query.Count()

		// Count should always equal number of results
		if c != n {
			t.Errorf("Testing %s. Number of results does not equal count. Count: %v, Results: %v", testCase.testName, c, n)
		}

		// Test number of records retrieved
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Test Count
		if c != testCase.expected {
			t.Errorf("Testing %s (count). Expected %v - got %v", testCase.testName, testCase.expected, c)
		}

	}

}

// Index Queries

// Helper for making groups of depatments
func getDept(i int) int {
	if i <= 10 {
		return 1
	} else if i <= 20 {
		return 2
	} else {
		return 3
	}
}

func Test_IndexQuery_ExactMatch_DateRange(t *testing.T) {
	var orders []Tormentable

	for i := 1; i <= 30; i++ {
		order := &Order{
			Model: Model{
				ID: gouuidv6.NewFromTime(time.Date(2009, time.November, i, 23, 0, 0, 0, time.UTC)),
			},
			Department: getDept(i),
		}

		orders = append(orders, order)
	}

	randomiseTormentables(orders)

	db, _ := OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName       string
		department     interface{}
		addFrom, addTo bool
		from, to       time.Time
		expected       int
	}{
		{"match department 1 - no date restriction", 1, false, false, time.Time{}, time.Time{}, 10},
		{"match department 1 - from beginning of time", 1, true, false, time.Time{}, time.Now(), 10},
		{"match department 1 - from beginning of time to now", 1, true, true, time.Time{}, time.Now(), 10},
		{"match department 1 - from now (no to)", 1, true, false, time.Now(), time.Time{}, 0},
		{"match department 1 - from 1st Nov (no to)", 1, true, false, time.Date(2009, time.November, 1, 23, 0, 0, 0, time.UTC), time.Time{}, 10},
		{"match department 1 - from 5th Nov", 1, true, false, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Time{}, 6},
		{"match department 1 - from 1st-5th Nov", 1, true, true, time.Date(2009, time.November, 1, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), 5},
	}

	for _, testCase := range testCases {
		rangequeryResults := []Order{}
		query := db.Find(&rangequeryResults).Where("department", testCase.department)
		if testCase.addFrom {
			query = query.From(testCase.from)
		}

		if testCase.addTo {
			query = query.To(testCase.to)
		}

		n, _ := query.Run()

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

	}
}

// Index searching
func Test_IndexQuery_Range(t *testing.T) {
	// Set up 100 orders with increasing department, customer and shipping fee
	// and save
	var orders []Tormentable

	for i := 0; i < 100; i++ {
		orders = append(orders, &Order{
			Department:  i + 1,
			Customer:    fmt.Sprintf("customer-%v", string(i+65)),
			ShippingFee: float64(i) + 0.99,
		})
	}

	// Randomise order before saving,
	// to ensure save order is not affecting retrieval
	// in some roundabout way
	randomiseTormentables(orders)

	db, _ := OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName   string
		indexName  string
		start, end interface{}
		expected   int
	}{
		// Non existent index
		{"non existent index", "notanindex", nil, nil, 0},

		// Int
		{"integer", "department", nil, nil, 100},
		{"integer - from 1", "department", 1, nil, 100},
		{"integer - from 2", "department", 2, nil, 99},
		{"integer - from 50", "department", 50, nil, 51},
		{"integer - 1 to 2", "department", 1, 2, 2},
		{"integer - 50 to 59", "department", 50, 59, 10},
		{"integer - 1 to 100", "department", 1, 100, 100},
		{"integer - to 50", "department", nil, 50, 50},

		// String
		{"string", "customer", nil, nil, 100},
		{"string", "customer", "customer", nil, 100},
		{"string - from A", "customer", "customer-A", nil, 100},
		{"string - from B", "customer", "customer-B", nil, 99},
		{"string - from Z", "customer", "customer-Z", nil, 75},
		{"string - from A to Z", "customer", "customer-A", "customer-Z", 26},
		{"string - to Z", "customer", nil, "customer-Z", 26},

		// Float
		{"float", "shippingfee", nil, nil, 100},
		{"float", "shippingfee", 0, nil, 100},
		{"float", "shippingfee", 0.99, nil, 100},
		{"float - from 1.99", "shippingfee", 1.99, nil, 99},
		{"float - from 50.99", "shippingfee", 50.99, nil, 50},
		{"float - from 99.99", "shippingfee", 99.99, nil, 1},
		{"float - to 20.99", "shippingfee", nil, 20.99, 21},
	}

	for _, testCase := range testCases {
		rangequeryResults := []Order{}
		n, _ := db.
			Find(&rangequeryResults).
			Where(testCase.indexName, testCase.start, testCase.end).
			Run()

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Check each member of the results for nil ID, customer and shipping fee
		for i, order := range rangequeryResults {
			if order.ID.IsNil() {
				t.Errorf("Testing %s.  Order no %v has nil ID", testCase.testName, i)
			}

			if order.Department == 0 {
				t.Errorf("Testing %s.  Order no %v has 0 department", testCase.testName, i)
			}

			if order.Customer == "" {
				t.Errorf("Testing %s.  Order no %v has blank customer", testCase.testName, i)
			}

			if order.ShippingFee == 0.0 {
				t.Errorf("Testing %s.  Order no %v has 0 shipping fee", testCase.testName, i)
			}
		}
	}

}

func Test_IndexQuery_Range_Plus_DateRange(t *testing.T) {
	var orders []Tormentable

	for i := 1; i <= 30; i++ {
		order := &Order{
			Department: getDept(i),
		}

		orders = append(orders, order)
	}

	randomiseTormentables(orders)

	db, _ := OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName   string
		start, end interface{}
		expected   int
	}{
		{"no range", nil, nil, 30},
		{"all departments", 1, 3, 30},
		{"departments 0, 1", 1, 2, 20},
		{"department 0", 1, 1, 10},
	}

	for _, testCase := range testCases {
		rangequeryResults := []Order{}
		n, _ := db.
			Find(&rangequeryResults).
			Where("department", testCase.start, testCase.end).
			Run()

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}
	}
}
