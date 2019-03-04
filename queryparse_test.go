package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Test_BuildQuery(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var results []testtypes.FullStruct

	testCases := []struct {
		testName    string
		queryString string
		targetQuery *tormenta.Query
		expectSame  bool
		expectError bool
	}{
		{
			"no query",
			"",
			db.Find(&results),
			true,
			false,
		},

		// Limit
		{
			"limit",
			"limit=1",
			db.Find(&results).Limit(1),
			true,
			false,
		},
		{
			"limit - different value",
			"limit=1",
			db.Find(&results).Limit(2),
			false,
			false,
		},
		{
			"limit - invalid value",
			"limit=word",
			db.Find(&results),
			true,
			true,
		},

		// Offset
		{
			"offset",
			"offset=1",
			db.Find(&results).Offset(1),
			true,
			false,
		},
		{
			"offset - different value",
			"offset=1",
			db.Find(&results).Offset(2),
			false,
			false,
		},
		{
			"offset - invalid value",
			"offset=word",
			db.Find(&results),
			true,
			true,
		},

		// Order
		{
			"order",
			"order=IntField",
			db.Find(&results).OrderBy("IntField"),
			true,
			false,
		},
		{
			"order - incorrect",
			"order=StringField",
			db.Find(&results).OrderBy("IntField"),
			false,
			false,
		},

		// Reverse
		{
			"reverse",
			"reverse=true",
			db.Find(&results).Reverse(),
			true,
			false,
		},
		{
			"reverse false should not match reversed query",
			"reverse=false",
			db.Find(&results).Reverse(),
			false,
			false,
		},
		{
			"reverse false should should match non reversed query",
			"reverse=false",
			db.Find(&results),
			true,
			false,
		},

		// Index searches
		{
			"index without anything else should error",
			"where=index:IntField",
			db.Find(&results),
			true,
			true,
		},
		{
			"match without index should error",
			"where=match:1",
			db.Find(&results),
			true,
			true,
		},
		{
			"start without index should error",
			"where=start:1",
			db.Find(&results),
			true,
			true,
		},
		{
			"end without index should error",
			"where=end:1",
			db.Find(&results),
			true,
			true,
		},
		{
			"match - correct",
			"where=index:IntField,match:1",
			db.Find(&results).Match("IntField", 1),
			true,
			false,
		},
		{
			"match - incorrect",
			"where=index:IntField,match:2",
			db.Find(&results).Match("IntField", 1),
			false,
			false,
		},

		// Range
		{
			"range - start only - correct",
			"where=index:IntField,start:1",
			db.Find(&results).Range("IntField", 1, nil),
			true,
			false,
		},
		{
			"range - end only - correct",
			"where=index:IntField,end:100",
			db.Find(&results).Range("IntField", nil, 100),
			true,
			false,
		},
		{
			"range - start and end - correct",
			"where=index:IntField,start:1,end:100",
			db.Find(&results).Range("IntField", 1, 100),
			true,
			false,
		},
		{
			"range - start and end - type mismatch",
			"where=index:IntField,start:1,end:invalidword",
			db.Find(&results),
			true,
			true,
		},
		{
			"index - match and range specified - no good",
			"where=index:IntField,start:1,end:100,match:1",
			db.Find(&results),
			true,
			true,
		},

		// From/To
		// Impossible to test equality because of use of UUID generation, but we can test for query building errors
		{
			"from - correct",
			"from=2006-01-02",
			db.Find(&results),
			false,
			false,
		},
		{
			"from - incorrect date format",
			"from=x-01-02",
			db.Find(&results),
			true,
			true,
		},
		{
			"to - correct",
			"to=2006-01-02",
			db.Find(&results),
			false,
			false,
		},
		{
			"to - incorrect date format",
			"to=x-01-02",
			db.Find(&results),
			true,
			true,
		},
		{
			"from is before to - possible",
			"from=2009-01-02&to=2010-01-02",
			db.Find(&results),
			false,
			false,
		},
		{
			"from is before to - possible - switch order in which they are specified",
			"to=2010-01-02&from=2009-01-02",
			db.Find(&results),
			false,
			false,
		},
		{
			"from is after to - impossible",
			"from=2010-01-02&to=2009-01-02",
			db.Find(&results),
			false, //does actually set the dates, but also errors
			true,
		},

		// Stack 'em up!
		{
			"limit, offset",
			"limit=1&offset=1",
			db.Find(&results).Limit(1).Offset(1),
			true,
			false,
		},
		{
			"limit, offset, reverse",
			"limit=1&offset=1&reverse=true",
			db.Find(&results).Limit(1).Offset(1).Reverse(),
			true,
			false,
		},
		{
			"limit, offset, reverse, index match",
			"limit=1&offset=1&reverse=true&where=index:IntField,match:1",
			db.Find(&results).Limit(1).Offset(1).Reverse().Match("IntField", 1),
			true,
			false,
		},
		{
			"limit, offset, reverse, index range",
			"limit=1&offset=1&reverse=true&where=index:IntField,start:1,end:10",
			db.Find(&results).Limit(1).Offset(1).Reverse().Range("IntField", 1, 10),
			true,
			false,
		},
		{
			"limit, offset, reverse, index range, startswith",
			"limit=1&offset=1&reverse=true&where=index:IntField,start:1,end:10&where=index:StringField,startswith:test",
			db.Find(&results).Limit(1).Offset(1).Reverse().Range("IntField", 1, 10).StartsWith("StringField", "test"),
			true,
			false,
		},
		{
			"limit, offset, reverse, index range, startswith, or",
			"or=true&limit=1&offset=1&reverse=true&where=index:IntField,start:1,end:10&where=index:StringField,startswith:test",
			db.Find(&results).Limit(1).Offset(1).Reverse().Range("IntField", 1, 10).Or().StartsWith("StringField", "test"),
			true,
			false,
		},
	}

	for _, test := range testCases {
		query := db.Find(&results)
		if err := query.Parse(false, test.queryString); err != nil && !test.expectError {
			t.Errorf("Testing %s. Building queries returned error: %s", test.testName, err)
		} else if test.expectError && err == nil {
			t.Errorf("Testing %s. Was expecting the built queries to error but it didn't", test.testName)
		}

		if test.expectSame {
			if !query.Compare(*test.targetQuery) {
				t.Errorf("Testing %s. Built query and target query are not equal. Expected %s, got %s", test.testName, *test.targetQuery, *query)
			}
		} else {
			if query.Compare(*test.targetQuery) {
				t.Errorf("Testing %s. Built query and target query are equal but I was expecting them to be different. Target: %s, built: %s", test.testName, *test.targetQuery, *query)
			}
		}

		// and finally just make sure that both the expected and built
		// queries actually run, catching errors in the test specification
		// We do this at the end, because hitting Run() casuses internal
		// change to the queries themselves, which messes up our equality checking
		if _, err := query.Run(); err != nil {
			t.Errorf("Testing %s. Built query returned error: %s", test.testName, err)
		}

		if _, err := test.targetQuery.Run(); err != nil {
			t.Errorf("Testing %s. Target query returned error: %s", test.testName, err)
		}
	}
}
