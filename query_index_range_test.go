package tormenta_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
)

// Test range queries across different types
func Test_IndexQuery_Range(t *testing.T) {
	// Set up 100 tts with increasing department, customer and shipping fee
	// and save
	var tts []tormenta.Record

	for i := 0; i < 100; i++ {
		tts = append(tts, &TestType{
			IntField:    i + 1,
			StringField: fmt.Sprintf("customer-%v", string((i%26)+65)),
			FloatField:  float64(i) + 0.99,
		})
	}

	// Randomise tt before saving,
	// to ensure save tt is not affecting retrieval
	// in some roundabout way
	tormenta.RandomiseRecords(tts)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(tts...)

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
		{"integer - no range", "intfield", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"integer - from 1", "intfield", 1, nil, 100, nil},
		{"integer - from 2", "intfield", 2, nil, 99, nil},
		{"integer - from 50", "intfield", 50, nil, 51, nil},
		{"integer - 1 to 2", "intfield", 1, 2, 2, nil},
		{"integer - 50 to 59", "intfield", 50, 59, 10, nil},
		{"integer - 1 to 100", "intfield", 1, 100, 100, nil},
		{"integer - to 50", "intfield", nil, 50, 50, nil},

		// String
		{"string - no range", "stringfield", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"string", "stringfield", "customer", nil, 100, nil},
		{"string - from A", "stringfield", "customer-A", nil, 100, nil},
		{"string - from B", "stringfield", "customer-B", nil, 96, nil},
		{"string - from Z", "stringfield", "customer-Z", nil, 3, nil},
		{"string - from A to Z", "stringfield", "customer-A", "customer-Z", 100, nil},
		{"string - to Z", "stringfield", nil, "customer-Z", 100, nil},

		// Float
		{"float - no range", "floatfield", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"float", "floatfield", 0, nil, 100, nil},
		{"float", "floatfield", 0.99, nil, 100, nil},
		{"float - from 1.99", "floatfield", 1.99, nil, 99, nil},
		{"float - from 50.99", "floatfield", 50.99, nil, 50, nil},
		{"float - from 99.99", "floatfield", 99.99, nil, 1, nil},
		{"float - to 20.99", "floatfield", nil, 20.99, 21, nil},
	}

	for _, testCase := range testCases {
		rangequeryResults := []TestType{}
		q := db.
			Find(&rangequeryResults).
			Range(testCase.indexName, testCase.start, testCase.end)

		// Forwards
		n, _, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number tts retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Check each member of the results for nil ID, customer and shipping fee
		for i, tt := range rangequeryResults {
			if tt.ID.IsNil() {
				t.Errorf("Testing %s.  Order no %v has nil ID", testCase.testName, i)
			}

			if tt.IntField == 0 {
				t.Errorf("Testing %s.  Order no %v has 0 department", testCase.testName, i)
			}

			if tt.StringField == "" {
				t.Errorf("Testing %s.  Order no %v has blank customer", testCase.testName, i)
			}

			if tt.FloatField == 0.0 {
				t.Errorf("Testing %s.  Order no %v has 0 shipping fee", testCase.testName, i)
			}
		}

		// Reverse
		rn, _, err := q.Reverse().Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number tts retrieved). Expected %v - got %v", testCase.testName, testCase.expected, rn)
		}

		// Check each member of the results for nil ID, customer and shipping fee
		for i, tt := range rangequeryResults {
			if tt.ID.IsNil() {
				t.Errorf("Testing %s.  Order no %v has nil ID", testCase.testName, i)
			}

			if tt.IntField == 0 {
				t.Errorf("Testing %s.  Order no %v has 0 department", testCase.testName, i)
			}

			if tt.StringField == "" {
				t.Errorf("Testing %s.  Order no %v has blank customer", testCase.testName, i)
			}

			if tt.FloatField == 0.0 {
				t.Errorf("Testing %s.  Order no %v has 0 shipping fee", testCase.testName, i)
			}
		}

	}

}

// Test index with multiple coinciding values
func Test_IndexQuery_Range_MultipleIndexMembers(t *testing.T) {
	var tts []tormenta.Record

	for i := 1; i <= 30; i++ {
		tt := &TestType{
			IntField: getDept(i),
		}

		tts = append(tts, tt)
	}

	tormenta.RandomiseRecords(tts)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(tts...)

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
		rangequeryResults := []TestType{}
		n, _, err := db.
			Find(&rangequeryResults).
			Range("intfield", testCase.start, testCase.end).
			Run()

		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number tts retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Reverse
		rangequeryResults = []TestType{}
		rn, _, err := db.
			Find(&rangequeryResults).
			Range("intfield", testCase.start, testCase.end).
			Reverse().
			Run()

		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number tts retrieved). Expected %v - got %v", testCase.testName, testCase.expected, rn)
		}
	}
}

// Test index queries augmented with a date range
func Test_IndexQuery_Range_DateRange(t *testing.T) {
	var tts []tormenta.Record

	for i := 1; i <= 30; i++ {
		tt := &TestType{
			Model: tormenta.Model{
				ID: gouuidv6.NewFromTime(time.Date(2009, time.November, i, 23, 0, 0, 0, time.UTC)),
			},
			IntField: getDept(i),
		}

		tts = append(tts, tt)
	}

	tormenta.RandomiseRecords(tts)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(tts...)

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
		rangequeryResults := []TestType{}
		query := db.Find(&rangequeryResults).Range("intfield", testCase.indexRangeStart, testCase.indexRangeEnd)

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
			t.Errorf("Testing %s (number tts retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Reverse
		nr, _, err := query.Reverse().Run()
		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number tts retrieved). Expected %v - got %v", testCase.testName, testCase.expected, nr)
		}

	}
}
