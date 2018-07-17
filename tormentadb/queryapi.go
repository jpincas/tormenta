package tormentadb

import (
	"errors"
	"time"

	"github.com/jpincas/gouuidv6"
)

// User API

type QueryOptions struct {
	First, Reverse bool
	Limit, Offset  int
	Start, End     interface{}
	From, To       time.Time
	IndexName      string
	IndexParams    []interface{}
}

// Find is the basic way to kick off a Query
func (db DB) Find(entities interface{}) *Query {
	return db.newQuery(entities, false)
}

// Query is another way of specifying a Query, using a struct of options instead of method chaining
func (db DB) Query(entities interface{}, options QueryOptions) *Query {
	q := db.newQuery(entities, options.First)

	// Overwrite limit if this is not a first-only search
	if !options.First {
		q.limit = options.Limit
	}

	if options.Offset > 0 {
		q.Offset(options.Offset)
	}

	// Apply reverse if speficied
	// Default is false, so can be left off
	q.reverse = options.Reverse

	// Apply date range if specified
	if !options.From.IsZero() {
		q.From(options.From)
	}

	if !options.To.IsZero() {
		q.To(options.To)
	}

	// Apply index if required
	if options.IndexName != "" {
		q.Where(options.IndexName, options.IndexParams...)
	}

	return q
}

// First kicks off a DB Query returning the first entity that matches the criteria
func (db DB) First(entity interface{}) *Query {
	return db.newQuery(entity, true)
}

// Limit limits the number of results a Query will return to n
func (q *Query) Limit(n int) *Query {
	q.limit = n
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

// Where takes an index name and up to 2 paramaters.
// If one parameter is supplied, the search is an exact match search
// If 2 parameters are supplied, it is a range search
func (q *Query) Where(indexName string, params ...interface{}) *Query {
	if len(params) == 1 {
		q.start = params[0]
		q.end = params[0]
	} else if len(params) == 2 {
		q.start = params[0]
		q.end = params[1]
	} else {
		q.err = errors.New("Index Query Where clause requires either 1 (exact match) or 2 (range) parameters")
	}

	q.isIndexQuery = true
	q.indexName = []byte(indexName)
	return q
}

// From adds a lower boundary to the date range of the Query
func (q *Query) From(t time.Time) *Query {
	// Subtract 1 nanosecond form the specified time
	// Leads to an inclusive date search
	t = t.Add(-1 * time.Nanosecond)

	q.from = gouuidv6.NewFromTime(t)
	return q
}

// To adds an upper bound to the date range of the Query
func (q *Query) To(t time.Time) *Query {

	q.to = gouuidv6.NewFromTime(t)
	return q
}

// Run actually executes the Query
func (q *Query) Run() (int, error) {
	return q.execute()
}

// Count executes the Query in fast, count-only mode
func (q *Query) Count() (int, error) {
	q.countOnly = true
	return q.execute()
}

func (q *Query) Sum(a interface{}) (int, error) {
	if !q.isIndexQuery || len(q.indexName) == 0 {
		return 0, errors.New("Aggregation must use an index Query")
	}

	q.aggTarget = a
	q.isAggQuery = true
	return q.execute()
}
