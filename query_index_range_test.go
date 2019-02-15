package tormenta_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

// Helpers

// just any old date really, all we're interested in is being able to sequence the year for testing
func dateWithYear(i int) time.Time {
	return time.Date(2009+i, time.November, 10, 23, 0, 0, 0, time.UTC)
}

func Test_IndexQuery_Range_Simple_Reverse(t *testing.T) {
	var fullStructs []tormenta.Record

	for i := 0; i < 100; i++ {
		fullStructs = append(fullStructs, &testtypes.FullStruct{
			IntField: i + 1,
		})
	}

	tormenta.RandomiseRecords(fullStructs)
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()
	db.Save(fullStructs...)

	rangequeryResults := []testtypes.FullStruct{}
	n, err := db.
		Find(&rangequeryResults).
		Range("intfield", 1, 100).
		Reverse().
		Run()

	if err != nil {
		t.Errorf("Testing simple reverse.  Got error %v", err)
	}

	if n != 100 {
		t.Errorf("Testing simple reverse.  Expected 100, got %v", n)
	}

}

// Test range queries across different types
func Test_IndexQuery_Range(t *testing.T) {
	// Set up 100 fullStructs with increasing department, customer and shipping fee
	// and save
	var fullStructs []tormenta.Record

	for i := 0; i < 100; i++ {
		fullStructs = append(fullStructs, &testtypes.FullStruct{
			IntField:          i + 1,
			DefinedIntField:   testtypes.DefinedInt(-50) + testtypes.DefinedInt(i),
			DefinedInt16Field: testtypes.DefinedInt16(-50) + testtypes.DefinedInt16(i),
			UintField:         uint(i) + 1,
			AnotherIntField:   (-50) + i,
			StringField:       fmt.Sprintf("customer-%v", string((i%26)+65)),
			FloatField:        float64(i) + 0.99,
			DefinedFloatField: testtypes.DefinedFloat(-50) + testtypes.DefinedFloat(i),
			AnotherFloatField: float64(-50) + float64(i),
			DateField:         dateWithYear(i),
		})
	}

	// Randomise fullStruct before saving,
	// to ensure save fullStruct is not affecting retrieval
	// in some roundabout way
	tormenta.RandomiseRecords(fullStructs)

	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	db.Save(fullStructs...)

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
		{"integer - start at 0", "intfield", 0, 100, 100, nil},
		{"integer - 1 to 100", "intfield", 1, 100, 100, nil},
		{"integer - to 50", "intfield", nil, 50, 50, nil},

		// Int (negatives - involves bit toggling - see interfaceToBytes())
		{"integer - start at -1, full range", "intfield", -1, 100, 100, nil},
		{"integer - start at -100, full range", "intfield", -100, 100, 100, nil},
		{"integer - start at -100, limited range ", "intfield", -100, -20, 0, nil},
		{"integer - negatives - out of range", "anotherintfield", -100, -51, 0, nil},
		{"integer - negatives - just in range", "anotherintfield", -100, -50, 1, nil},
		{"integer - negatives - first half of range", "anotherintfield", -50, -1, 50, nil},
		{"integer - negatives - span neg and pos - all", "anotherintfield", -50, 50, 100, nil},
		{"integer - negatives - span neg and pos", "anotherintfield", -25, 25, 51, nil},
		{"integer - negatives - span neg and pos 2", "anotherintfield", -10, 5, 16, nil},

		// Defined Int
		{"defined integer - no range", "definedintfield", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"defined integer - from 1", "definedintfield", 1, nil, 49, nil},
		{"defined integer - from 2", "definedintfield", 2, nil, 48, nil},
		{"defined integer - from 50", "definedintfield", 50, nil, 0, nil},
		{"defined integer - 1 to 2", "definedintfield", 1, 2, 2, nil},
		{"defined integer - 50 to 59", "definedintfield", 50, 59, 0, nil},
		{"defined integer - start at 0", "definedintfield", 0, 100, 50, nil},
		{"defined integer - 1 to 100", "definedintfield", 1, 100, 49, nil},
		{"defined integer - to 50", "definedintfield", nil, 50, 100, nil},
		{"defined integer - start at -1, full range", "definedintfield", -1, 100, 51, nil},
		{"defined integer - start at -100, full range", "definedintfield", -100, 100, 100, nil},
		{"defined integer - start at -100, limited range ", "definedintfield", -100, -20, 31, nil},
		{"defined integer - negatives - out of range", "definedintfield", -100, -51, 0, nil},
		{"defined integer - negatives - just in range", "definedintfield", -100, -50, 1, nil},
		{"defined integer - negatives - first half of range", "definedintfield", -50, -1, 50, nil},
		{"defined integer - negatives - span neg and pos - all", "definedintfield", -50, 50, 100, nil},
		{"defined integer - negatives - span neg and pos", "definedintfield", -25, 25, 51, nil},
		{"defined integer - negatives - span neg and pos 2", "definedintfield", -10, 5, 16, nil},

		// Defined Int16
		{"defined integer 16 - no range", "definedint16field", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"defined integer 16 - from 1", "definedint16field", int16(1), nil, 49, nil},
		{"defined integer 16 - from 2", "definedint16field", int16(2), nil, 48, nil},
		{"defined integer 16 - from 50", "definedint16field", int16(50), nil, 0, nil},
		{"defined integer 16 - 1 to 2", "definedint16field", int16(1), int16(2), 2, nil},
		{"defined integer 16 - 50 to 59", "definedint16field", int16(50), int16(59), 0, nil},
		{"defined integer 16 - start at 0", "definedint16field", int16(0), int16(100), 50, nil},
		{"defined integer 16 - 1 to 100", "definedint16field", int16(1), int16(100), 49, nil},
		{"defined integer 16 - to 50", "definedint16field", nil, int16(50), 100, nil},
		{"defined integer 16 - start at -1, full range", "definedint16field", int16(-1), int16(100), 51, nil},
		{"defined integer 16 - start at -100, full range", "definedint16field", int16(-100), int16(100), 100, nil},
		{"defined integer 16 - start at -100, limited range ", "definedint16field", int16(-100), int16(-20), 31, nil},
		{"defined integer 16 - negatives - out of range", "definedint16field", int16(-100), int16(-51), 0, nil},
		{"defined integer 16 - negatives - just in range", "definedint16field", int16(-100), int16(-50), 1, nil},
		{"defined integer 16 - negatives - first half of range", "definedint16field", int16(-50), int16(-1), 50, nil},
		{"defined integer 16 - negatives - span neg and pos - all", "definedint16field", int16(-50), int16(50), 100, nil},
		{"defined integer 16 - negatives - span neg and pos", "definedint16field", int16(-25), int16(25), 51, nil},
		{"defined integer 16 - negatives - span neg and pos 2", "definedint16field", int16(-10), int16(5), 16, nil},

		// Uint
		// Note how the types have to be explcitly stated, otherwise they will
		// be interpreted as ints and the sign bit will be flipped
		{"unsigned integer - no range", "uintfield", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"unsigned integer - from 1", "uintfield", uint(1), nil, 100, nil},
		{"unsigned integer - from 2", "uintfield", uint(2), nil, 99, nil},
		{"unsigned integer - from 50", "uintfield", uint(50), nil, 51, nil},
		{"unsigned integer - 1 to 2", "uintfield", uint(1), uint(2), 2, nil},
		{"unsigned integer - 50 to 59", "uintfield", uint(50), uint(59), 10, nil},
		{"unsigned integer - start at 0", "uintfield", uint(0), uint(100), 100, nil},
		{"unsigned integer - 1 to 100", "uintfield", uint(1), uint(100), 100, nil},
		{"unsigned integer - to 50", "uintfield", nil, uint(50), 50, nil},

		// Date - just encoded as an int64 so should be no problem
		{"date - no range", "datefield", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"date - all", "datefield", dateWithYear(0), nil, 100, nil},
		{"date - first 2", "datefield", dateWithYear(0), dateWithYear(1), 2, nil},
		{"date - random range", "datefield", dateWithYear(10), dateWithYear(20), 11, nil},

		// String
		{"string - no range", "stringfield", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"string", "stringfield", "customer", nil, 100, nil},
		{"string - from A", "stringfield", "customer-A", nil, 100, nil},
		{"string - from B", "stringfield", "customer-B", nil, 96, nil},
		{"string - from Z", "stringfield", "customer-Z", nil, 3, nil},
		{"string - from A to Z", "stringfield", "customer-A", "customer-Z", 100, nil},
		{"string - to Z", "stringfield", nil, "customer-Z", 100, nil},

		// Float
		// Note that we've always used the decimal point, so
		// the range values will be interpreted as floats not ints
		{"float - no range", "floatfield", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"float - 0 to nil", "floatfield", 0.00, nil, 100, nil},
		{"float - 0.99 to nil", "floatfield", 0.99, nil, 100, nil},
		{"float - from 1.99", "floatfield", 1.99, nil, 99, nil},
		{"float - from 50.99", "floatfield", 50.99, nil, 50, nil},
		{"float - from 99.99", "floatfield", 99.99, nil, 1, nil},
		{"float - to 20.99", "floatfield", nil, 20.99, 21, nil},

		// Negative floats
		// using regular floatField which has no negatives
		{"float - start at -1, full range", "floatfield", -1.00, 99.99, 100, nil},
		{"float - start at -100, full range", "floatfield", -100.00, 99.99, 100, nil},
		{"float - start at -100, limited range ", "floatfield", -100.00, -20.00, 0, nil},

		// now using anotherFloatField which does have negatives starting at -50
		{"float - negatives - just out of range", "anotherfloatfield", -100.00, -51.00, 0, nil},
		{"float - negatives - just in range", "anotherfloatfield", -100.00, -50.00, 1, nil},
		{"float - negatives - first half of range", "anotherfloatfield", -50.00, -1.00, 50, nil},
		{"float- negatives - span neg and pos", "anotherfloatfield", -50.00, 50.00, 100, nil},
		{"float- negatives - span neg and pos 2", "anotherfloatfield", -20.00, 30.00, 51, nil},

		// Defined Float
		{"defined float - no range", "definedfloatfield", nil, nil, 0, errors.New(tormenta.ErrNilInputsRangeIndexQuery)},
		{"defined float - 0 to nil", "definedfloatfield", 0.00, nil, 50, nil},
		{"defined float - 0.99 to nil", "definedfloatfield", 0.99, nil, 49, nil},
		{"defined float - from 1.99", "definedfloatfield", 1.99, nil, 48, nil},
		{"defined float - from 50.99", "definedfloatfield", 50.99, nil, 0, nil},
		{"defined float - from 99.99", "definedfloatfield", 99.99, nil, 0, nil},
		{"defined float - to 20.99", "definedfloatfield", nil, 20.99, 71, nil},
		{"defined float - start at -1, full range", "definedfloatfield", -1.00, 99.99, 51, nil},
		{"defined float - start at -100, full range", "definedfloatfield", -100.00, 99.99, 100, nil},
		{"defined float - start at -100, limited range ", "definedfloatfield", -100.00, -20.00, 31, nil},
		{"defined float - negatives - just out of range", "definedfloatfield", -100.00, -51.00, 0, nil},
		{"defined float - negatives - just in range", "definedfloatfield", -100.00, -50.00, 1, nil},
		{"defined float - negatives - first half of range", "definedfloatfield", -50.00, -1.00, 50, nil},
		{"defined float- negatives - span neg and pos", "definedfloatfield", -50.00, 50.00, 100, nil},
		{"defined float- negatives - span neg and pos 2", "definedfloatfield", -20.00, 30.00, 51, nil},
	}

	for _, testCase := range testCases {
		rangequeryResults := []testtypes.FullStruct{}
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
			t.Errorf("Testing %s (number fullStructs retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Check each member of the results for nil ID, customer and shipping fee
		for i, fullStruct := range rangequeryResults {
			if fullStruct.ID.IsNil() {
				t.Errorf("Testing %s.  Order no %v has nil ID", testCase.testName, i)
			}

			if fullStruct.IntField == 0 {
				t.Errorf("Testing %s.  Order no %v has 0 department", testCase.testName, i)
			}

			if fullStruct.StringField == "" {
				t.Errorf("Testing %s.  Order no %v has blank customer", testCase.testName, i)
			}

			if fullStruct.FloatField == 0.0 {
				t.Errorf("Testing %s.  Order no %v has 0 shipping fee", testCase.testName, i)
			}
		}

		// Reverse
		rangequeryResults = []testtypes.FullStruct{}
		q = db.
			Find(&rangequeryResults).
			Range(testCase.indexName, testCase.start, testCase.end).
			Reverse()
		rn, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s - reverse. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s - reverse. Didn't expect error [%v]", testCase.testName, err)
		}

		// Check for correct number of returned results
		if rn != testCase.expected || rn != n {
			t.Errorf("Testing %s - reverse (number fullStructs retrieved). Expected %v - got %v. Forwards search was %v", testCase.testName, testCase.expected, rn, n)
		}

		// Check each member of the results for nil ID, customer and shipping fee
		for i, fullStruct := range rangequeryResults {
			if fullStruct.ID.IsNil() {
				t.Errorf("Testing %s.  Order no %v has nil ID", testCase.testName, i)
			}

			if fullStruct.IntField == 0 {
				t.Errorf("Testing %s.  Order no %v has 0 department", testCase.testName, i)
			}

			if fullStruct.StringField == "" {
				t.Errorf("Testing %s.  Order no %v has blank customer", testCase.testName, i)
			}

			if fullStruct.FloatField == 0.0 {
				t.Errorf("Testing %s.  Order no %v has 0 shipping fee", testCase.testName, i)
			}
		}

	}

}

// Test index with multiple coinciding values
func Test_IndexQuery_Range_MultipleIndexMembers(t *testing.T) {
	var fullStructs []tormenta.Record

	for i := 1; i <= 30; i++ {
		fullStruct := &testtypes.FullStruct{
			IntField: getDept(i),
		}

		fullStructs = append(fullStructs, fullStruct)
	}

	tormenta.RandomiseRecords(fullStructs)

	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()
	db.Save(fullStructs...)

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
		rangequeryResults := []testtypes.FullStruct{}
		n, err := db.
			Find(&rangequeryResults).
			Range("intfield", testCase.start, testCase.end).
			Run()

		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number fullStructs retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Reverse
		rangequeryResults = []testtypes.FullStruct{}
		rn, err := db.
			Find(&rangequeryResults).
			Range("intfield", testCase.start, testCase.end).
			Reverse().
			Run()

		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number fullStructs retrieved). Expected %v - got %v", testCase.testName, testCase.expected, rn)
		}
	}
}

// Test index queries augmented with a date range
func Test_IndexQuery_Range_DateRange(t *testing.T) {
	var fullStructs []tormenta.Record

	for i := 1; i <= 30; i++ {
		fullStruct := &testtypes.FullStruct{
			Model: tormenta.Model{
				ID: gouuidv6.NewFromTime(time.Date(2009, time.November, i, 23, 0, 0, 0, time.UTC)),
			},
			IntField: getDept(i),
		}

		fullStructs = append(fullStructs, fullStruct)
	}

	tormenta.RandomiseRecords(fullStructs)

	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()
	db.Save(fullStructs...)

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
		rangequeryResults := []testtypes.FullStruct{}
		query := db.Find(&rangequeryResults).Range("intfield", testCase.indexRangeStart, testCase.indexRangeEnd)

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
			t.Errorf("Testing %s (number fullStructs retrieved). Expected %v - got %v", testCase.testName, testCase.expected, n)
		}

		// Reverse
		nr, err := query.Reverse().Run()
		if err != nil {
			t.Error("Testing basic querying - got error")
		}

		// Check for correct number of returned results
		if n != testCase.expected {
			t.Errorf("Testing %s (number fullStructs retrieved). Expected %v - got %v", testCase.testName, testCase.expected, nr)
		}

	}
}
