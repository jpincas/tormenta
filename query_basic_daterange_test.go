package tormenta_test

import (
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/demo"
)

func Test_BasicQuery_DateRange(t *testing.T) {
	// Create a list of orders over a date range
	var orders []tormenta.Tormentable
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
		orders = append(orders, &demo.Order{
			Model: tormenta.Model{
				ID: gouuidv6.NewFromTime(date),
			},
		})
	}

	// Save the orders
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	// Also another entity, to make sure there is no crosstalk
	db.Save(&demo.Product{
		Code:          "001",
		Name:          "Computer",
		Price:         999.99,
		StartingStock: 50,
		Description:   demo.DefaultDescription})

	//Quick check that all orders have saved correctly
	var results []demo.Order
	n, _, _ := db.Find(&results).Run()

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
		offset    int
	}{
		{"from right now - no orders expected, no 'to'", time.Now(), time.Time{}, 0, false, 0, 0},
		{"from beginning of time - all orders should be included, no 'to'", time.Time{}, time.Time{}, len(orders), false, 0, 0},
		{"from beginning of time - offset 1", time.Time{}, time.Time{}, len(orders) - 1, false, 0, 1},
		{"from beginning of time - offset 2", time.Time{}, time.Time{}, len(orders) - 2, false, 0, 2},
		{"from 2014, no 'to'", time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), time.Time{}, 8, false, 0, 0},
		{"from 1 hour ago, no 'to'", time.Now().Add(-1 * time.Hour), time.Time{}, 1, false, 0, 0},
		{"from beginning of time to now - expect all", time.Time{}, time.Now(), len(orders), true, 0, 0},
		{"from beginning of time to 2014 - expect 5", time.Time{}, time.Date(2014, time.January, 1, 1, 0, 0, 0, time.UTC), 5, true, 0, 0},
		{"from beginning of time to an hour ago - expect all but 1", time.Time{}, time.Now().Add(-1 * time.Hour), len(orders) - 1, true, 0, 0},
		{"from beginning of time - limit 1", time.Time{}, time.Time{}, 1, false, 1, 0},
		{"from beginning of time - limit 10", time.Time{}, time.Time{}, 10, false, 10, 0},
		{"from beginning of time - limit 10 - offset 2 (shouldnt affect number of results)", time.Time{}, time.Time{}, 10, false, 10, 2},
		{"from beginning of time - limit more than there are", time.Time{}, time.Time{}, len(orders), false, 0, 0},
	}

	for _, testCase := range testCases {
		rangequeryResults := []demo.Order{}
		query := db.Find(&rangequeryResults).From(testCase.from)

		if testCase.includeTo {
			query = query.To(testCase.to)
		}

		if testCase.limit > 0 {
			query = query.Limit(testCase.limit)
		}

		if testCase.offset > 0 {
			query = query.Offset(testCase.offset)
		}

		// FORWARD TESTS

		n, _, err := query.Run()
		if err != nil {
			t.Errorf("Testing %s. Got error %s", testCase.testName, err.Error())
		}

		c, _, err := query.Count()
		if err != nil {
			t.Errorf("Testing %s. Got error %s", testCase.testName, err.Error())
		}

		// Test number of records retrieved
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Test Count
		if c != testCase.expected {
			t.Errorf("Testing %s (count). Expected %v - got %v", testCase.testName, testCase.expected, c)
		}

		//Count should always equal number of results
		if c != n {
			t.Errorf("Testing %s. Number of results does not equal count. Count: %v, Results: %v", testCase.testName, c, n)
		}

		// REVERSE TESTS

		query = query.Reverse()

		rn, _, err := query.Run()
		if err != nil {
			t.Errorf("Testing REVERSE %s. Got error %s", testCase.testName, err.Error())
		}

		rc, _, err := query.Count()
		if err != nil {
			t.Errorf("Testing REVERSE %s. Got error %s", testCase.testName, err.Error())
		}

		// Test number of records retrieved
		if rn != testCase.expected {
			t.Errorf("Testing REVERSE %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, rn)
		}

		// Test Count
		if rc != testCase.expected {
			t.Errorf("Testing REVERSE %s (count). Expected %v - got %v", testCase.testName, testCase.expected, rc)
		}

		//Count should always equal number of results
		if rc != rn {
			t.Errorf("Testing REVERSE %s. Number of results does not equal count. Count: %v, Results: %v", testCase.testName, rc, rn)
		}

	}

}
