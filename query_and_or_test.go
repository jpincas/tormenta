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

	// Sum
	expectedOrSum  float64
	expectedAndSum float64
}

// Note the order in which we expect the results - date/time order!
func testCases(q *tormenta.Query) []orAndTest {
	return []orAndTest{
		{
			"single clause",
			[]*tormenta.Query{
				q.Cp().Match("intfield", 1),
			},
			1,
			[]testtypes.FullStruct{
				{IntField: 1},
			},
			1,
			[]testtypes.FullStruct{
				{IntField: 1},
			},
			1,
			1,
		},
		{
			"2 clauses",
			[]*tormenta.Query{
				q.Cp().Match("intfield", 1),
				q.Cp().Match("intfield", 2),
			},
			2,
			[]testtypes.FullStruct{
				{IntField: 1},
				{IntField: 2},
			},
			0,
			[]testtypes.FullStruct{},
			3,
			0,
		},
		{
			"more than 2 clauses",
			[]*tormenta.Query{
				q.Cp().Match("intfield", 1),
				q.Cp().Match("intfield", 2),
				q.Cp().Match("intfield", 3),
			},
			3,
			[]testtypes.FullStruct{
				{IntField: 1},
				{IntField: 2},
				{IntField: 3},
			},
			0,
			[]testtypes.FullStruct{},
			6,
			0,
		},
		{
			"more than 2 clauses - order of clauses should not matter",
			[]*tormenta.Query{
				q.Cp().Match("intfield", 2),
				q.Cp().Match("intfield", 1),
				q.Cp().Match("intfield", 3),
			},
			3,
			[]testtypes.FullStruct{
				{IntField: 1},
				{IntField: 2},
				{IntField: 3},
			},
			0,
			[]testtypes.FullStruct{},
			6,
			0,
		},
		{
			"more than 2 clauses - mixed indexes",
			[]*tormenta.Query{
				q.Cp().Match("intfield", 2),
				q.Cp().Match("stringfield", "int-1"),
				q.Cp().Match("intfield", 3),
			},
			3,
			[]testtypes.FullStruct{
				{IntField: 1},
				{IntField: 2},
				{IntField: 3},
			},
			0,
			[]testtypes.FullStruct{},
			6,
			0,
		},
		{
			"more than 2 clauses - mixed indexes - mixed matchers",
			[]*tormenta.Query{
				q.Cp().Range("intfield", 3, 5),
				q.Cp().Match("stringfield", "int-1"),
			},
			4,
			[]testtypes.FullStruct{
				{IntField: 1},
				{IntField: 3},
				{IntField: 4},
				{IntField: 5},
			},
			0,
			[]testtypes.FullStruct{},
			13,
			0,
		},
		{
			"more than 2 clauses - testing AND mainly",
			[]*tormenta.Query{
				q.Cp().Range("intfield", 1, 5),
				q.Cp().Match("stringfield", "int-2"),
			},
			5,
			[]testtypes.FullStruct{
				{IntField: 1},
				{IntField: 2},
				{IntField: 3},
				{IntField: 4},
				{IntField: 5},
			},
			1,
			[]testtypes.FullStruct{
				{IntField: 2},
			},
			15,
			2,
		},
		{
			"more than 2 clauses - testing AND in overlapping ranges",
			[]*tormenta.Query{
				q.Cp().Range("intfield", 1, 5),
				q.Cp().Range("stringfield", "int-2", "int-4"),
			},
			5,
			[]testtypes.FullStruct{
				{IntField: 1},
				{IntField: 2},
				{IntField: 3},
				{IntField: 4},
				{IntField: 5},
			},
			3,
			[]testtypes.FullStruct{
				{IntField: 2},
				{IntField: 3},
				{IntField: 4},
			},
			15,
			9,
		},
		{
			"nested OR",
			[]*tormenta.Query{
				tormenta.Or(
					q.Cp().Range("intfield", 1, 3),
					q.Cp().Range("stringfield", "int-1", "int-2"),
					// -> 1, 2, 3
				),
				tormenta.Or(
					q.Cp().Range("intfield", 4, 5),
					q.Cp().Range("stringfield", "int-5", "int-5"),
					// -> 4, 5
				),
			},
			5,
			[]testtypes.FullStruct{
				{IntField: 1},
				{IntField: 2},
				{IntField: 3},
				{IntField: 4},
				{IntField: 5},
			},
			0,
			[]testtypes.FullStruct{},
			15,
			0,
		},
		{
			"nested OR",
			[]*tormenta.Query{
				tormenta.Or(
					q.Cp().Range("intfield", 1, 4),
					q.Cp().Range("stringfield", "int-1", "int-2"),
					// -> 1, 2, 3, 4
				),
				tormenta.Or(
					q.Cp().Range("intfield", 4, 5),
					q.Cp().Range("stringfield", "int-5", "int-5"),
					// -> 4, 5
				),
			},
			5,
			[]testtypes.FullStruct{
				{IntField: 1},
				{IntField: 2},
				{IntField: 3},
				{IntField: 4},
				{IntField: 5},
			},
			1,
			[]testtypes.FullStruct{
				{IntField: 4},
			},
			15,
			4,
		},
		{
			"nested AND",
			[]*tormenta.Query{
				tormenta.And(
					q.Cp().Range("intfield", 1, 4),
					q.Cp().Range("stringfield", "int-1", "int-2"),
					// -> 1, 2
				),
				tormenta.And(
					q.Cp().Range("intfield", 4, 5),
					q.Cp().Range("stringfield", "int-5", "int-5"),
					// -> 5
				),
			},
			3,
			[]testtypes.FullStruct{
				{IntField: 1},
				{IntField: 2},
				{IntField: 5},
			},
			0,
			[]testtypes.FullStruct{},
			8,
			0,
		},
		{
			"nested AND",
			[]*tormenta.Query{
				tormenta.And(
					q.Cp().Range("intfield", 2, 4),
					q.Cp().Range("stringfield", "int-1", "int-4"),
					// -> 2, 3, 4
				),
				tormenta.Or(
					q.Cp().Range("intfield", 4, 5),
					q.Cp().Range("stringfield", "int-5", "int-5"),
					// -> 4, 5
				),
			},
			4,
			[]testtypes.FullStruct{
				{IntField: 2},
				{IntField: 3},
				{IntField: 4},
				{IntField: 5},
			},
			1,
			[]testtypes.FullStruct{
				{IntField: 4},
			},
			14,
			4,
		},
	}
}

func Test_And_Basic(t *testing.T) {
	// DB
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
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
	q := db.Find(&results)
	testCases := testCases(q)

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

		// Test 'Sum'

		sum, _, err := tormenta.And(testCase.clauses...).Sum([]string{"IntField"})

		if err != nil {
			t.Errorf("Testing basic AND (%s, sum) - got error: %v", testCase.testName, err)
		}

		if sum != testCase.expectedAndSum {
			t.Errorf("Testing basic AND (%s, sum). Wrong sum result. Expected: %v; got: %v", testCase.testName, testCase.expectedAndSum, sum)
		}
	}
}

// func Test_Or_Basic(t *testing.T) {
// 	// DB
// 	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
// 	defer db.Close()

// 	// Generate some simple data and save
// 	var toSave []tormenta.Record
// 	for i := 0; i < 10; i++ {
// 		toSave = append(toSave, &testtypes.FullStruct{
// 			IntField:    i,
// 			StringField: fmt.Sprintf("int-%v", i),
// 		})
// 	}
// 	db.Save(toSave...)

// 	// Results placeholder and generate test cases
// 	results := []testtypes.FullStruct{}
// 	testCases := testCases(db, &results)

// 	for _, testCase := range testCases {
// 		results = []testtypes.FullStruct{}

// 		////////
// 		// OR //
// 		////////

// 		// Test 'Run'

// 		n, err := tormenta.Or(testCase.clauses...).Run()

// 		if err != nil {
// 			t.Errorf("Testing basic OR (%s,run)- got error", testCase.testName)
// 		}

// 		if n != len(results) {
// 			t.Errorf("Testing basic OR (%s,run) - n does not equal length of results. N: %v; Length results: %v", testCase.testName, n, len(results))
// 		}

// 		if n != testCase.expectedOrN {
// 			t.Errorf("Testing basic OR (%s,run). Wrong number of results. Expected: %v; got: %v", testCase.testName, testCase.expectedOrN, n)
// 		}

// 		for i, _ := range results {
// 			if results[i].IntField != testCase.expectedOrResults[i].IntField {
// 				t.Errorf("Testing basic OR (%s,run). Mismatch in array member %v", testCase.testName, i)
// 			}
// 		}

// 		// Test 'Count'

// 		c, err := tormenta.Or(testCase.clauses...).Count()

// 		if err != nil {
// 			t.Errorf("Testing basic OR (%s,count) - got error", testCase.testName)
// 		}

// 		if c != testCase.expectedOrN {
// 			t.Errorf("Testing basic OR (%s,count). Wrong number of results. Expected: %v; got: %v", testCase.testName, testCase.expectedOrN, c)
// 		}

// 		// Test 'Sum'

// 		sum, _, err := tormenta.Or(testCase.clauses...).Sum([]string{"IntField"})

// 		if err != nil {
// 			t.Errorf("Testing basic OR (%s, sum) - got error: %v", testCase.testName, err)
// 		}

// 		if sum != testCase.expectedOrSum {
// 			t.Errorf("Testing basic OR (%s, sum). Wrong sum result. Expected: %v; got: %v", testCase.testName, testCase.expectedOrSum, sum)
// 		}

// 	}
// }
