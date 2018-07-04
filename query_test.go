package tormenta

import (
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
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

	c, err := db.Query(&orders).Count()
	if c != len(orders) {
		t.Errorf("Testing count 1 entity saved. Expecting 1 - got %v", c)
	}

	// 2 orders
	order2 := Order{}
	orders = []Order{}
	db.Save(&order2)

	n, _ = db.Query(&orders).Run()

	if len(orders) != 2 || n != 2 {
		t.Errorf("Testing querying with 2 entity saved. Expecting 2 entities - got %v/%v", len(orders), n)
	}

	c, err = db.Query(&orders).Count()
	if c != len(orders) {
		t.Errorf("Testing count 2 entities saved. Expecting 2 - got %v", c)
	}

}

func Test_RangeQuery(t *testing.T) {
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
	n, _ := db.Query(&results).Run()

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
		rangeQueryResults := []Order{}
		query := db.Query(&rangeQueryResults).From(testCase.from)

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
