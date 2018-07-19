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

func Test_IndexQuery_StartsWith(t *testing.T) {
	customers := []string{"j", "jo", "jon", "jonathan", "job", "pablo"}
	var orders []tormenta.Tormentable

	for _, customer := range customers {
		orders = append(orders, &demo.Order{
			Customer: customer,
		})
	}

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName      string
		startsWith    string
		reverse       bool
		expected      int
		expectedError error
	}{
		{"blank string", "", false, 0, errors.New(tormenta.ErrBlankInputStartsWithQuery)},
		{"no match - no interference", "nocustomerwiththisname", false, 0, nil},
		{"single match - no interference", "pablo", false, 1, nil},
		{"single match - possible interference", "jonathan", false, 1, nil},
		{"single match - possible interference", "job", false, 1, nil},
		{"wide match - 1 letter", "j", false, 5, nil},
		{"wide match - 2 letters", "jo", false, 4, nil},
		{"wide match - 3 letters", "jon", false, 2, nil},

		// Reversed - shouldn't make any difference to N
		{"blank string", "", true, 0, errors.New(tormenta.ErrBlankInputStartsWithQuery)},
		{"no match - no interference", "nocustomerwiththisname", true, 0, nil},
		{"single match - no interference", "pablo", true, 1, nil},
		{"single match - possible interference", "jonathan", true, 1, nil},
		{"single match - possible interference", "job", true, 1, nil},
		{"wide match - 1 letter", "j", true, 5, nil},
		{"wide match - 2 letters", "jo", true, 4, nil},
		{"wide match - 3 letters", "jon", true, 2, nil},
	}

	for _, testCase := range testCases {
		results := []demo.Order{}

		q := db.Find(&results).StartsWith("customer", testCase.startsWith)
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
