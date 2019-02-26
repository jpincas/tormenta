package urltoquery

import (
	"net/url"
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Test_BuildQuery(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var results []testtypes.FullStruct

	testCases := []struct {
		testName      string
		queryString   string
		targetQueries []*tormenta.Query
		expectSame    bool
		expectError   bool
	}{
		{
			"no query",
			"",
			[]*tormenta.Query{},
			true,
			false,
		},

		// Limit
		{
			"limit",
			"query=limit:1",
			[]*tormenta.Query{db.Find(&results).Limit(1)},
			true,
			false,
		},
		{
			"limit - different value",
			"query=limit:1",
			[]*tormenta.Query{db.Find(&results).Limit(2)},
			false,
			false,
		},
		{
			"limit - invalid value",
			"query=limit:word",
			[]*tormenta.Query{},
			false,
			true,
		},

		// Offset
		{
			"offset",
			"query=offset:1",
			[]*tormenta.Query{db.Find(&results).Offset(1)},
			true,
			false,
		},
		{
			"offset - different value",
			"query=offset:1",
			[]*tormenta.Query{db.Find(&results).Offset(2)},
			false,
			false,
		},
		{
			"offset - invalid value",
			"query=offset:word",
			[]*tormenta.Query{},
			false,
			true,
		},

		// Reverse
		{
			"reverse",
			"query=reverse:true",
			[]*tormenta.Query{db.Find(&results).Reverse()},
			true,
			false,
		},
		{
			"reverse false should not match reversed query",
			"query=reverse:false",
			[]*tormenta.Query{db.Find(&results).Reverse()},
			false,
			false,
		},
		{
			"reverse false should should match non reversed query",
			"query=reverse:false",
			[]*tormenta.Query{db.Find(&results)},
			true,
			false,
		},

		// // Index searches
		{
			"index without anything else should error",
			"query=index:indexstring",
			[]*tormenta.Query{},
			false,
			true,
		},
		{
			"match without index should error",
			"query=match:1",
			[]*tormenta.Query{},
			false,
			true,
		},
		{
			"start without index should error",
			"query=start:1",
			[]*tormenta.Query{},
			false,
			true,
		},
		{
			"end without index should error",
			"query=end:1",
			[]*tormenta.Query{},
			false,
			true,
		},
		{
			"match - correct",
			"query=index:indexstring,match:1",
			[]*tormenta.Query{db.Find(&results).Match("indexstring", 1)},
			true,
			false,
		},
		{
			"match - incorrect",
			"query=index:indexstring,match:2",
			[]*tormenta.Query{db.Find(&results).Match("indexstring", 1)},
			false,
			false,
		},

		// Range
		{
			"range - start only - correct",
			"query=index:indexstring,start:1",
			[]*tormenta.Query{db.Find(&results).Range("indexstring", 1, nil)},
			true,
			false,
		},
		{
			"range - end only - correct",
			"query=index:indexstring,end:100",
			[]*tormenta.Query{db.Find(&results).Range("indexstring", nil, 100)},
			true,
			false,
		},
		{
			"range - start and end - correct",
			"query=index:indexstring,start:1,end:100",
			[]*tormenta.Query{db.Find(&results).Range("indexstring", 1, 100)},
			true,
			false,
		},
		{
			"range - start and end - type mismatch",
			"query=index:indexstring,start:1,end:invalidword",
			[]*tormenta.Query{},
			false,
			true,
		},
		{
			"index - match and range specified - no good",
			"query=index:indexstring,start:1,end:100,match:1",
			[]*tormenta.Query{},
			false,
			true,
		},

		// From/To
		// Impossible to test equality because of use of UUID generation, but we can test for query building errors
		{
			"from - correct",
			"query=from:2006-01-02",
			[]*tormenta.Query{db.Find(&results)},
			false,
			false,
		},
		{
			"from - incorrect date format",
			"query=from:x-01-02",
			[]*tormenta.Query{},
			true,
			true,
		},
		{
			"to - correct",
			"query=to:2006-01-02",
			[]*tormenta.Query{db.Find(&results)},
			false,
			false,
		},
		{
			"to - incorrect date format",
			"query=to:x-01-02",
			[]*tormenta.Query{},
			true,
			true,
		},
		{
			"from is before to - possible",
			"query=from:2009-01-02,to:2010-01-02",
			[]*tormenta.Query{db.Find(&results)},
			false,
			false,
		},
		{
			"from is before to - possible - switch order in which they are specified",
			"query=to:2010-01-02,from:2009-01-02",
			[]*tormenta.Query{db.Find(&results)},
			false,
			false,
		},
		{
			"from is after to - impossible",
			"query=from:2010-01-02,to:2009-01-02",
			[]*tormenta.Query{},
			false,
			true,
		},

		// Stack 'em up!
		{
			"limit, offset",
			"query=limit:1,offset:1",
			[]*tormenta.Query{db.Find(&results).Limit(1).Offset(1)},
			true,
			false,
		},
		{
			"limit, offset, reverse",
			"query=limit:1,offset:1,reverse:true",
			[]*tormenta.Query{db.Find(&results).Limit(1).Offset(1).Reverse()},
			true,
			false,
		},
		{
			"limit, offset, reverse, index match",
			"query=limit:1,offset:1,reverse:true,index:indexstring,match:1",
			[]*tormenta.Query{db.Find(&results).Limit(1).Offset(1).Reverse().Match("indexstring", 1)},
			true,
			false,
		},
		{
			"limit, offset, reverse, index range",
			"query=limit:1,offset:1,reverse:true,index:indexstring,start:1,end:10",
			[]*tormenta.Query{db.Find(&results).Limit(1).Offset(1).Reverse().Range("indexstring", 1, 10)},
			true,
			false,
		},

		// Multiple queries
		{
			"multiple queries - both ok",
			"query=limit:1,offset:1&query=limit:1,offset:2",
			[]*tormenta.Query{
				db.Find(&results).Limit(1).Offset(1),
				db.Find(&results).Limit(1).Offset(2),
			},
			true,
			false,
		},
		{
			"multiple queries - 1 bad, expect overall error",
			"query=limit:rubbishword,offset:1&query=limit:1,offset:2",
			[]*tormenta.Query{},
			false,
			true,
		},
	}

	for _, test := range testCases {
		// Parse the query string for values
		values, err := url.ParseQuery(test.queryString)
		if err != nil {
			t.Errorf("Testing %s. Parsing URL returned error: %s", test.testName, err)
		}

		// Attempt to build the queries,
		// but don't combine them as it would make it impossible
		// to test that each has been built properly
		queries, err := buildQueries(db, &results, values, false)
		if err != nil && !test.expectError {
			t.Errorf("Testing %s. Building queries returned error: %s", test.testName, err)
		} else if test.expectError && err == nil {
			t.Errorf("Testing %s. Was expecting the built queries to error but it didn't", test.testName)
		}

		if len(queries) != len(test.targetQueries) {
			t.Fatalf("Testing %s. Was expecting %v queries but got %v", test.testName, len(test.targetQueries), len(queries))
		}

		// Finally, iterate through the built queries and
		// actually make sure the two queries are equal if thats the aim
		for i, q := range queries {
			if test.expectSame {
				if !q.Compare(*test.targetQueries[i]) {
					t.Errorf("Testing %s - index %v. Built query and target query are not equal. Expected %s, got %s", test.testName, i, *test.targetQueries[i], *q)
				}
			} else {
				if q.Compare(*test.targetQueries[i]) {
					t.Errorf("Testing %s - index %v. Built query and target query are equal but I was expecting them to be different. Target: %s, built: %s", test.testName, i, *test.targetQueries[i], *q)
				}
			}

			// and finally just make sure that both the expected and built
			// queries actually run, catching errors in the test specification
			// We do this at the end, because hitting Run() casuses internal
			// change to the queries themselves, which messes up our equality checking
			if _, err := q.Run(); err != nil {
				t.Errorf("Testing %s. Built query returned error: %s", test.testName, err)
			}

			if _, err := test.targetQueries[i].Run(); err != nil {
				t.Errorf("Testing %s. Target query returned error: %s", test.testName, err)
			}
		}

	}

}
