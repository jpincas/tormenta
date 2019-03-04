package tormenta

import (
	"errors"
	"strings"
	"time"

	"github.com/jpincas/gouuidv6"
)

// QUERY INITIATORS

// Find is the basic way to kick off a Query
func (db DB) Find(entities interface{}) *Query {
	return db.newQuery(entities)
}

// First kicks off a DB Query returning the first entity that matches the criteria
func (db DB) First(entity interface{}) *Query {
	q := db.newQuery(entity)
	q.limit = 1
	q.single = true
	return q
}

// Debug turns on helpful debugging information for the query
func (q *Query) Debug() *Query {
	q.debug = true
	return q
}

// CONTEXT SETTING

// SetContext allows a context to be passed through the query
func (q *Query) SetContext(key string, val interface{}) *Query {
	q.ctx[key] = val
	return q
}

// FILTER APPLICATION

// Match adds an exact-match index search to a query
func (q *Query) Match(indexName string, param interface{}) *Query {
	// For a single parameter 'exact match' search, it is non sensical to pass nil
	// Set the error and return the query unchanged
	if param == nil {
		q.err = errors.New(ErrNilInputMatchIndexQuery)
		return q
	}

	// If we are matching a string, lower-case it
	switch param.(type) {
	case string:
		param = strings.ToLower(param.(string))
	}

	// Create the filter and add it on
	q.addFilter(filter{
		start:     param,
		end:       param,
		indexName: toIndexName(indexName),
	})

	return q
}

// Range adds a range-match index search to a query
func (q *Query) Range(indexName string, start, end interface{}) *Query {
	// For an index range search,
	// it is non-sensical to pass two nils
	// Set the error and return the query unchanged
	if start == nil && end == nil {
		q.err = errors.New(ErrNilInputsRangeIndexQuery)
		return q
	}

	// Create the filter and add it on
	q.addFilter(filter{
		start:     start,
		end:       end,
		indexName: toIndexName(indexName),
	})

	return q
}

// StartsWith allows for string prefix filtering
func (q *Query) StartsWith(indexName string, s string) *Query {
	// Blank string is not valid
	if s == "" {
		q.err = errors.New(ErrBlankInputStartsWithQuery)
		return q
	}

	// Create the filter and add it on
	q.addFilter(filter{
		start:             s,
		end:               s,
		isStartsWithQuery: true,
		indexName:         toIndexName(indexName),
	})

	return q
}

// GLOBAL QUERY MODIFIERS

// Sets the query to return filter results combined in a logical OR way instead of AND.
// It doesn't matter where in the chain, you put it - all filters will be combined in an OR
// fashion if it appears just once.  Having said that, if you are combining two filters, it
// reads nicely to put the Or() in the middle, e.g.
// .Range("myint", 1, 10).Or().StartsWith("mystring", "test"),
func (q *Query) Or() *Query {
	q.idsCombinator = union
	return q
}

// Sets the query to return filter results combined in a logical AND way.  This is the default,
// so this should rarely be necessary.  Mainly useful for the query parser.
func (q *Query) And() *Query {
	q.idsCombinator = intersection
	return q
}

// Limit limits the number of results a Query will return to n.
// If a limit has already been set on a query and you try to set a new one, it will only
// be overriden if it is lower.  This allows you easily set a 'hard' limit up front,
// that cannot be overriden for that query.
func (q *Query) Limit(n int) *Query {
	if q.limit == 0 {
		q.limit = n
	} else if n < q.limit {
		q.limit = n
	}

	return q
}

// Offset starts N entities from the beginning
func (q *Query) Offset(n int) *Query {
	q.offset = n
	q.offsetCounter = n
	return q
}

// Reverse reverses the order of date range scanning and returned results (i.e. scans from 'new' to 'old', instead of the default 'old' to 'new' )
func (q *Query) Reverse() *Query {
	q.reverse = true
	return q
}

// UnReverse unsets reverse on a query. Not expected to be particularly useful but needed by the string to query builder
func (q *Query) UnReverse() *Query {
	q.reverse = false
	return q
}

// OrderBy specifies an index by which to order results..
func (q *Query) OrderBy(indexName string) *Query {
	q.orderByIndexName = toIndexName(indexName)
	return q
}

// From adds a lower boundary to the date range of the Query
func (q *Query) From(t time.Time) *Query {
	q.from = fromUUID(t)
	return q
}

// To adds an upper bound to the date range of the Query
func (q *Query) To(t time.Time) *Query {
	q.to = toUUID(t)
	return q
}

// ManualFromToSet allows you to set the exact gouuidv6s for from and to
// Useful for testing purposes.
func (q *Query) ManualFromToSet(from, to gouuidv6.UUID) *Query {
	q.from = from
	q.to = to
	return q
}

// QUERY EXECUTORS

// Run actually executes the Query
func (q *Query) Run() (int, error) {
	return q.execute()
}

// Count executes the Query in fast, count-only mode
func (q *Query) Count() (int, error) {
	q.countOnly = true
	return q.execute()
}

// Sum produces a sum aggregation using the index only, which is much faster
// than accessing every record
func (q *Query) Sum(a interface{}, indexName string) (int, error) {
	q.sumTarget = a
	q.sumIndexName = toIndexName(indexName)
	return q.execute()
}
