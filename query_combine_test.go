package tormenta_test

import (
	"fmt"
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

type orAndTest struct {
	testName string
	clauses  []*tormenta.Query

	// For 'OR' tests
	expectedOrN       int
	expectedOrResults []testtypes.FullStruct

	// For 'AND' tests
	expectedAndN       int
	expectedAndResults []testtypes.FullStruct
}

// Note the order in which we expect the results - date/time order!
func testCases(db *tormenta.DB, results *[]testtypes.FullStruct) []orAndTest {
	return []orAndTest{
		{
			"single clause",
			[]*tormenta.Query{
				db.Find(results).Match("intfield", 1),
			},
			1,
			[]testtypes.FullStruct{
				{IntField: 1},
			},
			1,
			[]testtypes.FullStruct{
				{IntField: 1},
			},
		},
		{
			"2 clauses",
			[]*tormenta.Query{
				db.Find(results).Match("intfield", 1),
				db.Find(results).Match("intfield", 2),
			},
			2,
			[]testtypes.FullStruct{
				{IntField: 2},
				{IntField: 1},
			},
			0,
			[]testtypes.FullStruct{},
		},
		{
			"more than 2 clauses",
			[]*tormenta.Query{
				db.Find(results).Match("intfield", 1),
				db.Find(results).Match("intfield", 2),
				db.Find(results).Match("intfield", 3),
			},
			3,
			[]testtypes.FullStruct{
				{IntField: 3},
				{IntField: 2},
				{IntField: 1},
			},
			0,
			[]testtypes.FullStruct{},
		},
		{
			"more than 2 clauses - order of clauses should not matter",
			[]*tormenta.Query{
				db.Find(results).Match("intfield", 2),
				db.Find(results).Match("intfield", 1),
				db.Find(results).Match("intfield", 3),
			},
			3,
			[]testtypes.FullStruct{
				{IntField: 3},
				{IntField: 2},
				{IntField: 1},
			},
			0,
			[]testtypes.FullStruct{},
		},
		{
			"more than 2 clauses - mixed indexes",
			[]*tormenta.Query{
				db.Find(results).Match("intfield", 2),
				db.Find(results).Match("stringfield", "int-1"),
				db.Find(results).Match("intfield", 3),
			},
			3,
			[]testtypes.FullStruct{
				{IntField: 3},
				{IntField: 2},
				{IntField: 1},
			},
			0,
			[]testtypes.FullStruct{},
		},
		{
			"more than 2 clauses - mixed indexes - mixed matchers",
			[]*tormenta.Query{
				db.Find(results).Range("intfield", 3, 5),
				db.Find(results).Match("stringfield", "int-1"),
			},
			4,
			[]testtypes.FullStruct{
				{IntField: 5},
				{IntField: 4},
				{IntField: 3},
				{IntField: 1},
			},
			0,
			[]testtypes.FullStruct{},
		},
		{
			"more than 2 clauses - testing AND mainly",
			[]*tormenta.Query{
				db.Find(results).Range("intfield", 1, 5),
				db.Find(results).Match("stringfield", "int-2"),
			},
			5,
			[]testtypes.FullStruct{
				{IntField: 5},
				{IntField: 4},
				{IntField: 3},
				{IntField: 2},
				{IntField: 1},
			},
			1,
			[]testtypes.FullStruct{
				{IntField: 2},
			},
		},
		{
			"more than 2 clauses - testing AND in overlapping ranges",
			[]*tormenta.Query{
				db.Find(results).Range("intfield", 1, 5),
				db.Find(results).Range("stringfield", "int-2", "int-4"),
			},
			5,
			[]testtypes.FullStruct{
				{IntField: 5},
				{IntField: 4},
				{IntField: 3},
				{IntField: 2},
				{IntField: 1},
			},
			3,
			[]testtypes.FullStruct{
				{IntField: 4},
				{IntField: 3},
				{IntField: 2},
			},
		},
	}
}

func Test_And_Basic(t *testing.T) {
	// DB
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	// Generate some simple data and save
	var toSave []tormenta.Record
	for i := 0; i < 10; i++ {
		toSave = append(toSave, &testtypes.FullStruct{
			IntField:    i,
			StringField: fmt.Sprintf("int-%v", i),
		})
	}
	db.Save(toSave...)

	// Results placeholder and generate test cases
	results := []testtypes.FullStruct{}
	testCases := testCases(db, &results)

	for _, testCase := range testCases {
		results = []testtypes.FullStruct{}

		// Test 'Run'

		n, err := tormenta.And(testCase.clauses...).Run()

		if err != nil {
			t.Errorf("Testing basic AND (%s,run)- got error", testCase.testName)
		}

		if n != len(results) {
			t.Errorf("Testing basic AND (%s,run) - n does not equal length of results. N: %v; Length results: %v", testCase.testName, n, len(results))
		}

		if n != testCase.expectedAndN {
			t.Errorf("Testing basic AND (%s,run). Wrong number of results. Expected: %v; got: %v", testCase.testName, testCase.expectedAndN, n)
		}

		for i, _ := range results {
			if results[i].IntField != testCase.expectedAndResults[i].IntField {
				t.Errorf("Testing basic AND (%s,run). Mismatch in array member %v", testCase.testName, i)
			}
		}

		// Test 'Count'

		c, err := tormenta.And(testCase.clauses...).Count()

		if err != nil {
			t.Errorf("Testing basic AND (%s,count) - got error", testCase.testName)
		}

		if c != testCase.expectedAndN {
			t.Errorf("Testing basic AND (%s,count). Wrong number of results. Expected: %v; got: %v", testCase.testName, testCase.expectedAndN, c)
		}

	}
}

func Test_Or_Basic(t *testing.T) {
	// DB
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	// Generate some simple data and save
	var toSave []tormenta.Record
	for i := 0; i < 10; i++ {
		toSave = append(toSave, &testtypes.FullStruct{
			IntField:    i,
			StringField: fmt.Sprintf("int-%v", i),
		})
	}
	db.Save(toSave...)

	// Results placeholder and generate test cases
	results := []testtypes.FullStruct{}
	testCases := testCases(db, &results)

	for _, testCase := range testCases {
		results = []testtypes.FullStruct{}

		////////
		// OR //
		////////

		// Test 'Run'

		n, err := tormenta.Or(testCase.clauses...).Run()

		if err != nil {
			t.Errorf("Testing basic OR (%s,run)- got error", testCase.testName)
		}

		if n != len(results) {
			t.Errorf("Testing basic OR (%s,run) - n does not equal length of results. N: %v; Length results: %v", testCase.testName, n, len(results))
		}

		if n != testCase.expectedOrN {
			t.Errorf("Testing basic OR (%s,run). Wrong number of results. Expected: %v; got: %v", testCase.testName, testCase.expectedOrN, n)
		}

		for i, _ := range results {
			if results[i].IntField != testCase.expectedOrResults[i].IntField {
				t.Errorf("Testing basic OR (%s,run). Mismatch in array member %v", testCase.testName, i)
			}
		}

		// Test 'Count'

		c, err := tormenta.Or(testCase.clauses...).Count()

		if err != nil {
			t.Errorf("Testing basic OR (%s,count) - got error", testCase.testName)
		}

		if c != testCase.expectedOrN {
			t.Errorf("Testing basic OR (%s,count). Wrong number of results. Expected: %v; got: %v", testCase.testName, testCase.expectedOrN, c)
		}

	}
}
