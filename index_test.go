package tormenta

import (
	"fmt"
	"testing"

	"github.com/dgraph-io/badger"
)

// Index Creation
func Test_BasicIndexing(t *testing.T) {
	db, _ := OpenTest("data/tests")
	defer db.Close()

	// Create basic order and save
	// Orders have an 'index' on Customer field
	order := Order{
		Department:  99,
		Customer:    "jon",
		ShippingFee: 5.99,
	}

	db.Save(&order)

	db.KV.View(func(txn *badger.Txn) error {
		customerIndex := makeIndexKey([]byte("order"), order.ID, "customer", "jon")
		departmentIndex := makeIndexKey([]byte("order"), order.ID, "department", 99)
		shippingFeeIndex := makeIndexKey([]byte("order"), order.ID, "shippingfee", 5.99)

		_, err := txn.Get(customerIndex)
		if err == badger.ErrKeyNotFound {
			t.Error("Testing basic index key setting (string). Could not get index key")
		}

		_, err = txn.Get(departmentIndex)
		if err == badger.ErrKeyNotFound {
			t.Error("Testing basic index key setting (int). Could not get index key")
		}

		_, err = txn.Get(shippingFeeIndex)
		if err == badger.ErrKeyNotFound {
			t.Error("Testing basic index key setting (float64). Could not get index key")
		}

		return nil
	})
}

// Index searching
func Test_IndexRange(t *testing.T) {
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

	db, _ := OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	// For now we are only testing the number of returned results,
	// Not the actual returned values
	testCases := []struct {
		testName  string
		indexName string
		from, to  interface{}
		expected  int
	}{
		// Non existent index
		{"non existent index", "notanindex", nil, nil, 0},

		// Int
		{"integer", "department", nil, nil, 100},
		{"integer - from 1", "department", 1, nil, 100},
		{"integer - from 2", "department", 2, nil, 99},
		{"integer - from 50", "department", 50, nil, 51},
		{"integer - 50 to 59", "department", 50, 59, 10},
		{"integer - 1 to 101", "department", 0, 100, 100},
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
		n, _ := db.FindIndex(&rangequeryResults, testCase.indexName).From(testCase.from).To(testCase.to).Run()

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
