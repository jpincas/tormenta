package tormentadb_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// Test range queries across different types
func Test_IndexQuery_Range(t *testing.T) {
	// Set up 100 orders with increasing department, customer and shipping fee
	// and save
	var orders []tormenta.Tormentable

	for i := 0; i < 100; i++ {
		orders = append(orders, &demo.Order{
			Department:  i + 1,
			Customer:    fmt.Sprintf("customer-%v", string((i%26)+65)),
			ShippingFee: float64(i) + 0.99,
		})
	}

	// Randomise order before saving,
	// to ensure save order is not affecting retrieval
	// in some roundabout way
	tormenta.RandomiseTormentables(orders)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName      string
		indexName     string
		start, end    interface{}
		expected      int
		expectedError error
	}{
		// FORWARD

		// Non existent index
		{"non existent index - no range", "notanindex", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"non existent index", "notanindex", 1, 2, 0, nil},

		// Int
		{"integer - no range", "department", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"integer - from 1", "department", 1, nil, 100, nil},
		{"integer - from 2", "department", 2, nil, 99, nil},
		{"integer - from 50", "department", 50, nil, 51, nil},
		{"integer - 1 to 2", "department", 1, 2, 2, nil},
		{"integer - 50 to 59", "department", 50, 59, 10, nil},
		{"integer - 1 to 100", "department", 1, 100, 100, nil},
		{"integer - to 50", "department", nil, 50, 50, nil},

		// String
		{"string - no range", "customer", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"string", "customer", "customer", nil, 100, nil},
		{"string - from A", "customer", "customer-A", nil, 100, nil},
		{"string - from B", "customer", "customer-B", nil, 96, nil},
		{"string - from Z", "customer", "customer-Z", nil, 3, nil},
		{"string - from A to Z", "customer", "customer-A", "customer-Z", 100, nil},
		{"string - to Z", "customer", nil, "customer-Z", 100, nil},

		// Float
		{"float - no range", "shippingfee", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"float", "shippingfee", 0, nil, 100, nil},
		{"float", "shippingfee", 0.99, nil, 100, nil},
		{"float - from 1.99", "shippingfee", 1.99, nil, 99, nil},
		{"float - from 50.99", "shippingfee", 50.99, nil, 50, nil},
		{"float - from 99.99", "shippingfee", 99.99, nil, 1, nil},
		{"float - to 20.99", "shippingfee", nil, 20.99, 21, nil},
	}

	for _, testCase := range testCases {
		rangequeryResults := []demo.Order{}
		q := db.
			Find(&rangequeryResults).
			Range(testCase.indexName, testCase.start, testCase.end)

		// Forwards
		n, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

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

		// Reverse
		rn, err := q.Reverse().Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, rn)
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

// Test index with multiple coinciding values
func Test_IndexQuery_Range_MultipleIndexMembers(t *testing.T) {
	var orders []tormenta.Tormentable

	for i := 1; i <= 30; i++ {
		order := &demo.Order{
			Department: getDept(i),
		}

		orders = append(orders, order)
	}

	tormenta.RandomiseTormentables(orders)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName   string
		start, end interface{}
		expected   int
	}{
		{"all departments", 1, 3, 30},
		{"departments 1, 2", 1, 2, 20},
		{"department 1", 1, 1, 10},
	}

	for _, testCase := range testCases {
		// Forwards
		rangequeryResults := []demo.Order{}
		n, err := db.
			Find(&rangequeryResults).
			Range("department", testCase.start, testCase.end).
			Run()

		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Reverse
		rangequeryResults = []demo.Order{}
		rn, err := db.
			Find(&rangequeryResults).
			Range("department", testCase.start, testCase.end).
			Reverse().
			Run()

		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, rn)
		}
	}
}

// Test index queries augmented with a date range
func Test_IndexQuery_Range_DateRange(t *testing.T) {
	var orders []tormenta.Tormentable

	for i := 1; i <= 30; i++ {
		order := &demo.Order{
			Model: tormenta.Model{
				ID: gouuidv6.NewFromTime(time.Date(2009, time.November, i, 23, 0, 0, 0, time.UTC)),
			},
			Department: getDept(i),
		}

		orders = append(orders, order)
	}

	tormenta.RandomiseTormentables(orders)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName        string
		indexRangeStart interface{}
		addFrom, addTo  bool
		from, to        time.Time
		expected        int
		indexRangeEnd   interface{}
	}{
		// Range match tests
		{"departments 1-3 - no date restriction", 1, false, false, time.Time{}, time.Time{}, 30, 3},
		{"departments 1-3 - from beginning of time", 1, true, false, time.Time{}, time.Time{}, 30, 3},
		{"departments 1-3 - from beginning of time to now", 1, true, true, time.Time{}, time.Now(), 30, 3},
		{"departments 1-3 - from now (no to)", 1, true, false, time.Now(), time.Time{}, 0, 3},
		{"departments 1-3 - from 5th Nov", 1, true, false, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Time{}, 26, 3},
		{"departments 1-3 - from 5th Nov - 15th Nov", 1, true, true, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 15, 23, 0, 0, 0, time.UTC), 11, 3},
		{"departments 1-2 - from 5th Nov - 15th Nov", 1, true, true, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 15, 23, 0, 0, 0, time.UTC), 11, 2},
		{"departments 1-2 - from 5th Nov - 9th Nov", 1, true, true, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 9, 23, 0, 0, 0, time.UTC), 5, 2},
	}

	for _, testCase := range testCases {
		rangequeryResults := []demo.Order{}
		query := db.Find(&rangequeryResults).Range("department", testCase.indexRangeStart, testCase.indexRangeEnd)

		if testCase.addFrom {
			query = query.From(testCase.from)
		}

		if testCase.addTo {
			query = query.To(testCase.to)
		}

		// Forwards
		n, err := query.Run()
		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Reverse
		nr, err := query.Reverse().Run()
		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, nr)
		}

	}
}
