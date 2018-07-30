package tormentadb_test

import (
	"errors"
	"fmt"
	"testing"

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
		reverse       bool
	}{
		// FORWARD

		// Non existent index
		{"non existent index - no range", "notanindex", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery), false},
		{"non existent index", "notanindex", 1, 2, 0, nil, false},

		// Int
		{"integer - no range", "department", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery), false},
		{"integer - from 1", "department", 1, nil, 100, nil, false},
		{"integer - from 2", "department", 2, nil, 99, nil, false},
		{"integer - from 50", "department", 50, nil, 51, nil, false},
		{"integer - 1 to 2", "department", 1, 2, 2, nil, false},
		{"integer - 50 to 59", "department", 50, 59, 10, nil, false},
		{"integer - 1 to 100", "department", 1, 100, 100, nil, false},
		{"integer - to 50", "department", nil, 50, 50, nil, false},

		// String
		{"string - no range", "customer", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery), false},
		{"string", "customer", "customer", nil, 100, nil, false},
		{"string - from A", "customer", "customer-A", nil, 100, nil, false},
		{"string - from B", "customer", "customer-B", nil, 96, nil, false},
		{"string - from Z", "customer", "customer-Z", nil, 3, nil, false},
		{"string - from A to Z", "customer", "customer-A", "customer-Z", 100, nil, false},
		{"string - to Z", "customer", nil, "customer-Z", 100, nil, false},

		// Float
		{"float - no range", "shippingfee", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery), false},
		{"float", "shippingfee", 0, nil, 100, nil, false},
		{"float", "shippingfee", 0.99, nil, 100, nil, false},
		{"float - from 1.99", "shippingfee", 1.99, nil, 99, nil, false},
		{"float - from 50.99", "shippingfee", 50.99, nil, 50, nil, false},
		{"float - from 99.99", "shippingfee", 99.99, nil, 1, nil, false},
		{"float - to 20.99", "shippingfee", nil, 20.99, 21, nil, false},

		// REVERSE

		// Non existent index
		{"non existent index - no range", "notanindex", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery), true},
		{"non existent index", "notanindex", 1, 2, 0, nil, true},

		// Int
		{"integer - no range", "department", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery), true},
		{"integer - from 100 - reverse", "department", 100, nil, 100, nil, true},
		{"integer - from 99", "department", 99, nil, 99, nil, true},
		{"integer - from 50", "department", 50, nil, 50, nil, true},
		{"integer - 2 to 1", "department", 2, 1, 2, nil, true},
		{"integer - 59 to 50", "department", 59, 50, 10, nil, true},
		{"integer - 100 to 1", "department", 100, 1, 100, nil, true},
		{"integer - to 50", "department", nil, 50, 51, nil, true},

		// String
		{"string - no range", "customer", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery), true},
		{"string", "customer", "customer", nil, 100, nil, true},
		{"string - from A", "customer", "customer-A", nil, 4, nil, true},
		{"string - from B", "customer", "customer-B", nil, 8, nil, true},
		{"string - from Z", "customer", "customer-Z", nil, 100, nil,
			true},
		{"string - from Z to A", "customer", "customer-Z", "customer-A", 100, nil, true},
		{"string - to Z", "customer", nil, "customer-Z", 3, nil, true},

		// Float
		{"float - no range", "shippingfee", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery), true},
		{"float - from 0", "shippingfee", 0, nil, 0, nil, true},
		{"float - from 0.99", "shippingfee", 0.99, nil, 1, nil, true},
		{"float - from 1.99", "shippingfee", 1.99, nil, 2, nil, true},
		{"float - from 50.99", "shippingfee", 50.99, nil, 51, nil, true},
		{"float - from 99.99", "shippingfee", 99.99, nil, 100, nil, true},
		{"float - to 20.99", "shippingfee", nil, 20.99, 80, nil, true},
	}

	for _, testCase := range testCases {
		rangequeryResults := []demo.Order{}
		q := db.
			Find(&rangequeryResults).
			Range(testCase.indexName, testCase.start, testCase.end)

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
