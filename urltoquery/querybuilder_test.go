package urltoquery

import (
	"net/url"
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Test_BuildQuery(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	var results []testtypes.FullStruct

	// Queries default to/from endpoints using uuids created by time
	// so they will usually be different for each query.
	// In order to create identical queries that we can later compare,
	// we need a function that can delivery queries with identical from and to endpoints
	from := gouuidv6.NewFromTime(time.Date(1980, time.January, 1, 1, 0, 0, 0, time.UTC))
	to := gouuidv6.NewFromTime(time.Date(2080, time.January, 1, 1, 0, 0, 0, time.UTC))
	blankQuery := func() *tormenta.Query {
		q := db.Find(&results)
		q.ManualFromToSet(from, to)
		return q
	}

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
			blankQuery(),
			true,
			false,
		},

		// Limit
		{
			"limit",
			"limit=1",
			blankQuery().Limit(1),
			true,
			false,
		},
		{
			"limit - different value",
			"limit=1",
			blankQuery().Limit(2),
			false,
			false,
		},
		{
			"limit - invalid value",
			"limit=word",
			blankQuery().Limit(1),
			false,
			true,
		},

		// Offset
		{
			"offset",
			"offset=1",
			blankQuery().Offset(1),
			true,
			false,
		},
		{
			"offset - different value",
			"offset=1",
			blankQuery().Offset(2),
			false,
			false,
		},
		{
			"offset - invalid value",
			"offset=word",
			blankQuery().Offset(1),
			false,
			true,
		},

		// Reverse
		{
			"reverse",
			"reverse=true",
			blankQuery().Reverse(),
			true,
			false,
		},
		{
			"reverse false should not match reversed query",
			"reverse=false",
			blankQuery().Reverse(),
			false,
			false,
		},
		{
			"reverse false should should match non reversed query",
			"reverse=false",
			blankQuery(),
			true,
			false,
		},

		// Index searches
		{
			"index without anything else should error",
			"index=indexstring",
			blankQuery().Match("indexstring", 1),
			false,
			true,
		},
		{
			"match without index should error",
			"match=1",
			blankQuery().Match("indexstring", 1),
			false,
			true,
		},
		{
			"start without index should error",
			"start=1",
			blankQuery().Match("indexstring", 1),
			false,
			true,
		},
		{
			"end without index should error",
			"end=1",
			blankQuery().Match("indexstring", 1),
			false,
			true,
		},
		{
			"match - correct",
			"index=indexstring&match=1",
			blankQuery().Match("indexstring", 1),
			true,
			false,
		},
		{
			"range - start only - correct",
			"index=indexstring&start=1",
			blankQuery().Range("indexstring", 1, nil),
			true,
			false,
		},
		{
			"range - end only - correct",
			"index=indexstring&end=100",
			blankQuery().Range("indexstring", nil, 100),
			true,
			false,
		},
		{
			"range - start and end - correct",
			"index=indexstring&start=1&end=100",
			blankQuery().Range("indexstring", 1, 100),
			true,
			false,
		},
		{
			"range - start and end - type mismatch",
			"index=indexstring&start=1&end=invalidword",
			blankQuery().Range("indexstring", 1, 100),
			false,
			true,
		},
		{
			"index - match and range specified",
			"index=indexstring&start=1&end=100&match=1",
			blankQuery().Range("indexstring", 1, 100),
			false,
			true,
		},

		// From/To
		// Impossible to test equality because of use of UUID generation, but we can test for query building errors
		{
			"from - correct",
			"from=2006-01-02",
			blankQuery(),
			false,
			false,
		},
		{
			"from - incorrect date format",
			"from=x-01-02",
			blankQuery(),
			true,
			true,
		},
		{
			"to - correct",
			"to=2006-01-02",
			blankQuery(),
			false,
			false,
		},
		{
			"to - incorrect date format",
			"to=x-01-02",
			blankQuery(),
			true,
			true,
		},
		{
			"from is before to - possible",
			"from=2009-01-02&to=2010-01-02",
			blankQuery(),
			false,
			false,
		},
		{
			"from is before to - possible - switch order in which they are specified",
			"to=2010-01-02&from=2009-01-02",
			blankQuery(),
			false,
			false,
		},
		{
			"from is after to - impossible",
			"from=2010-01-02&to=2009-01-02",
			blankQuery(),
			false,
			true,
		},

		// Stack 'em up!
		{
			"limit, offset",
			"limit=1&offset=1",
			blankQuery().Limit(1).Offset(1),
			true,
			false,
		},
		{
			"limit, offset, reverse",
			"limit=1&offset=1&reverse=true",
			blankQuery().Limit(1).Offset(1).Reverse(),
			true,
			false,
		},
		{
			"limit, offset, reverse, index match",
			"limit=1&offset=1&reverse=true&index=indexstring&match=1",
			blankQuery().Limit(1).Offset(1).Reverse().Match("indexstring", 1),
			true,
			false,
		},
		{
			"limit, offset, reverse, index range",
			"limit=1&offset=1&reverse=true&index=indexstring&start=1&end=10",
			blankQuery().Limit(1).Offset(1).Reverse().Range("indexstring", 1, 10),
			true,
			false,
		},
	}

	for _, test := range testCases {
		// First test that the expected query doesn't error
		if _, err := test.targetQuery.Run(); err != nil {
			t.Errorf("Testing %s. Target query returned error: %s", test.testName, err)
		}

		// Parse the query string for values
		values, err := url.ParseQuery(test.queryString)
		if err != nil {
			t.Errorf("Testing %s. Parsing URL returned error: %s", test.testName, err)
		}

		// Set up a base query and then build it out
		q := blankQuery()
		if err := BuildQuery(q, values, false); err != nil && !test.expectError {
			t.Errorf("Testing %s. Building query returned error: %s", test.testName, err)
		} else if test.expectError && err == nil {
			t.Errorf("Testing %s. Was expecting the built query to error but it didn't", test.testName)
		}

		// Run it and make sure it doesn't error
		if _, err := q.Run(); err != nil {
			t.Errorf("Testing %s. Built query returned error: %s", test.testName, err)
		}

		// Finally, actually make sure the two queries are equal if thats the aim
		if test.expectSame {
			if !q.Compare(*test.targetQuery) {
				t.Errorf("Testing %s. Built query and target query are not equal. Expected %s, got %s", test.testName, *test.targetQuery, *q)
			}
		} else {
			if q.Compare(*test.targetQuery) {
				t.Errorf("Testing %s. Built query and target query are equal but I was expecting them to be different. Target: %s, built: %s", test.testName, *test.targetQuery, *q)
			}
		}

	}

}
