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

// Simple test of bool indexing
func Test_IndexQuery_Bool(t *testing.T) {
	orderFalse := demo.Order{}
	orderTrue := demo.Order{HasShipped: true}
	orderTrue2 := demo.Order{HasShipped: true}
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(&orderFalse, &orderTrue, &orderTrue2)

	results := []demo.Order{}
	// Test true
	n, _ := db.Find(&results).Match("hasshipped", true).Run()
	if n != 2 {
		t.Errorf("Testing bool index.  Expected 2 results, got %v", n)
	}

	// Test false + count
	c, _ := db.Find(&results).Match("hasshipped", false).Count()
	if c != 1 {
		t.Errorf("Testing bool index.  Expected 1 result, got %v", c)
	}
}

// Test exact matching on strings
func Test_IndexQuery_Match_String(t *testing.T) {
	customers := []string{"jon", "jonathan", "pablo"}
	var orders []tormenta.Tormentable

	for i := 0; i < 100; i++ {
		orders = append(orders, &demo.Order{
			Customer: customers[i%len(customers)],
		})
	}

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName      string
		match         interface{}
		reverse       bool
		expected      int
		expectedError error
	}{
		{"blank string", nil, false, 0, errors.New(tormenta.ErrNilInputMatchIndexQuery)},
		{"blank string", "", false, 0, nil},
		{"should not match any", "nocustomerwiththisname", false, 0, nil},
		{"matches 1 exactly with no interference", "pablo", false, 33, nil},
		{"matches 1 exactly and 1 prefix", "jon", false, 34, nil},
		{"matches 1 exactly and has same prefix as other", "jonathan", false, 33, nil},

		// Reversed - shouldn't make any difference to N
		{"blank string - reversed", nil, true, 0, errors.New(tormenta.ErrNilInputMatchIndexQuery)},
		{"blank string - reversed", "", true, 0, nil},
		{"should not match any - reversed", "nocustomerwiththisname", true, 0, nil},
		{"matches 1 exactly with no interference - reversed", "pablo", true, 33, nil},
		{"matches 1 exactly and 1 prefix - reversed", "jon", true, 34, nil},
		{"matches 1 exactly and has same prefix as other - reversed", "jonathan", true, 33, nil},
	}

	for _, testCase := range testCases {
		results := []demo.Order{}

		q := db.Find(&results).Match("customer", testCase.match)
		if testCase.reverse {
			q.Reverse()
		}

		n, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		if n != testCase.expected {
			t.Errorf("Testing %s.  Expecting %v, got %v", testCase.testName, testCase.expected, n)
		}
	}
}

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
		n, err := db.
			Find(&rangequeryResults).
			Range(testCase.indexName, testCase.start, testCase.end).
			Run()

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
	}

}

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
		rangequeryResults := []demo.Order{}
		n, _ := db.
			Find(&rangequeryResults).
			Range("department", testCase.start, testCase.end).
			Run()

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}
	}
}

// Test aggregation on an index
func Test_Aggregation(t *testing.T) {
	var products []tormenta.Tormentable

	for i := 1; i <= 30; i++ {
		product := &demo.Product{
			Price:         float64(i),
			StartingStock: i,
		}

		products = append(products, product)
	}

	tormenta.RandomiseTormentables(products)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(products...)

	results := []demo.Product{}
	var intSum int32
	var floatSum float64
	expected := 465

	// Int32

	_, err := db.Find(&results).Range("startingstock", 1, 30).Sum(&intSum)
	if err != nil {
		t.Error("Testing int32 agreggation.  Got error")
	}

	expectedIntSum := int32(expected)
	if intSum != expectedIntSum {
		t.Errorf("Testing int32 agreggation. Expteced %v, got %v", expectedIntSum, intSum)
	}

	// Float64

	_, err = db.Find(&results).Range("price", 1.00, 30.00).Sum(&floatSum)
	if err != nil {
		t.Error("Testing float64 agreggation.  Got error")
	}

	expectedFloatSum := float64(expected)
	if floatSum != expectedFloatSum {
		t.Errorf("Testing float64 agreggation. Expteced %v, got %v", expectedFloatSum, floatSum)
	}
}

// Test index queries augmented with a date range
func Test_IndexQuery_DateRange(t *testing.T) {
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
		// Exact match tests (indexRangeEnd is nil)
		{"match department 1 - no date restriction", 1, false, false, time.Time{}, time.Time{}, 10, nil},
		{"match department 1 - from beginning of time", 1, true, false, time.Time{}, time.Now(), 10, nil},
		{"match department 1 - from beginning of time to now", 1, true, true, time.Time{}, time.Now(), 10, nil},
		{"match department 1 - from now (no to)", 1, true, false, time.Now(), time.Time{}, 0, nil},
		{"match department 1 - from 1st Nov (no to)", 1, true, false, time.Date(2009, time.November, 1, 23, 0, 0, 0, time.UTC), time.Time{}, 10, nil},
		{"match department 1 - from 5th Nov", 1, true, false, time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), time.Time{}, 6, nil},
		{"match department 1 - from 1st-5th Nov", 1, true, true, time.Date(2009, time.November, 1, 23, 0, 0, 0, time.UTC), time.Date(2009, time.November, 5, 23, 0, 0, 0, time.UTC), 5, nil},

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

		// If only indexRangeStart is specified then its an exact match search
		query := db.Find(&rangequeryResults)
		if testCase.indexRangeEnd == nil {
			query = query.Match("department", testCase.indexRangeStart)
		} else {
			query = query.Range("department", testCase.indexRangeStart, testCase.indexRangeEnd)
		}

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
