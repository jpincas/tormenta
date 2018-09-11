package tormentadb_test

import (
	"errors"
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// Simple test of bool indexing
func Test_IndexQuery_Match_Bool(t *testing.T) {
	orderFalse := demo.Order{}
	orderTrue := demo.Order{HasShipped: true}
	orderTrue2 := demo.Order{HasShipped: true}
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(&orderFalse, &orderTrue, &orderTrue2)

	results := []demo.Order{}
	// Test true
	n, _, err := db.Find(&results).Match("hasshipped", true).Run()
	if err != nil {
		t.Error("Testing basic querying - got error")
	}

	if n != 2 {
		t.Errorf("Testing bool index.  Expected 2 results, got %v", n)
	}

	// Test false + count
	c, _, err := db.Find(&results).Match("hasshipped", false).Count()
	if err != nil {
		t.Error("Testing basic querying - got error")
	}

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
		testName string
		match    interface{}

		expected      int
		expectedError error
	}{
		{"blank string", nil, 0, errors.New(tormenta.ErrNilInputMatchIndexQuery)},
		{"blank string", "", 0, nil},
		{"should not match any", "nocustomerwiththisname", 0, nil},
		{"matches 1 exactly with no interference", "pablo", 33, nil},
		{"matches 1 exactly and 1 prefix", "jon", 34, nil},
		{"matches 1 exactly and has same prefix as other", "jonathan", 33, nil},
		{"uppercase - should make no difference", "JON", 34, nil},
		{"mixed-case - should make no difference", "Jon", 34, nil},
	}

	for _, testCase := range testCases {
		results := []demo.Order{}

		// Forwards
		q := db.Find(&results).Match("customer", testCase.match)
		n, _, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		if n != testCase.expected {
			t.Errorf("Testing %s.  Expecting %v, got %v", testCase.testName, testCase.expected, n)
		}

		// Reverse
		q = db.Find(&results).Match("customer", testCase.match).Reverse()
		rn, _, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		if n != testCase.expected {
			t.Errorf("Testing %s.  Expecting %v, got %v", testCase.testName, testCase.expected, rn)
		}
	}
}

func Test_IndexQuery_Match_Int(t *testing.T) {
	var orders []tormenta.Tormentable

	for i := 0; i < 100; i++ {
		orders = append(orders, &demo.Order{
			Department: i % 10,
		})
	}

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName      string
		match         interface{}
		expected      int
		expectedError error
	}{
		{"nothing", nil, 0, errors.New(tormenta.ErrNilInputMatchIndexQuery)},
		{"1", 1, 10, nil},
		{"11", 11, 0, nil},
	}

	for _, testCase := range testCases {
		results := []demo.Order{}

		// Forwards
		q := db.Find(&results).Match("department", testCase.match)
		n, _, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		if n != testCase.expected {
			t.Errorf("Testing %s.  Expecting %v, got %v", testCase.testName, testCase.expected, n)
		}

		// Reverse
		q = db.Find(&results).Match("customer", testCase.match).Reverse()
		rn, _, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		if n != testCase.expected {
			t.Errorf("Testing %s.  Expecting %v, got %v", testCase.testName, testCase.expected, rn)
		}
	}
}

func Test_IndexQuery_Match_Float(t *testing.T) {
	var orders []tormenta.Tormentable

	for i := 1; i <= 100; i++ {
		orders = append(orders, &demo.Order{
			ShippingFee: float64(i) / float64(10),
		})
	}

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName      string
		match         interface{}
		expected      int
		expectedError error
	}{
		{"nothing", nil, 0, errors.New(tormenta.ErrNilInputMatchIndexQuery)},
		{"0.1", 0.1, 1, nil},
		{"0.1", 0.10, 1, nil},
		{"0.11", 0.1, 1, nil},
		{"0.20", 0.200, 1, nil},
	}

	for _, testCase := range testCases {
		results := []demo.Order{}

		// Forwards
		q := db.Find(&results).Match("shippingfee", testCase.match)
		n, _, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		if n != testCase.expected {
			t.Errorf("Testing %s.  Expecting %v, got %v", testCase.testName, testCase.expected, n)
		}

		// Reverse
		q = db.Find(&results).Match("customer", testCase.match).Reverse()
		rn, _, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		if n != testCase.expected {
			t.Errorf("Testing %s.  Expecting %v, got %v", testCase.testName, testCase.expected, rn)
		}
	}
}
func Test_IndexQuery_Match_DateRange(t *testing.T) {
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
	}

	for _, testCase := range testCases {
		rangequeryResults := []demo.Order{}
		query := db.Find(&rangequeryResults).Match("department", testCase.indexRangeStart)

		if testCase.addFrom {
			query = query.From(testCase.from)
		}

		if testCase.addTo {
			query = query.To(testCase.to)
		}

		// Forwards
		n, _, err := query.Run()
		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		if n != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Backwards
		rn, _, err := query.Reverse().Run()
		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		if rn != testCase.expected {
			t.Errorf("Testing %s (number orders retrieved). Expected %v - got %v", testCase.testName, testCase.expected, rn)
		}

	}
}
