package archive

import "time"

type QueryOptions struct {
	First, Reverse bool
	Limit, Offset  int
	Start, End     interface{}
	From, To       time.Time
	IndexName      string
	IndexParams    []interface{}
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
	// Use 'match' for 1 param, 'range' for 2
	if options.IndexName != "" {
		if len(options.IndexParams) == 1 {
			q.Match(options.IndexName, options.IndexParams[0])
		} else if len(options.IndexParams) == 2 {
			q.Range(options.IndexName, options.IndexParams[0], options.IndexParams[1])
		}
	}

	return q
}

// Sum takes a slightly sifferent approach to aggregation - you might call it 'slow sum'.
// It doesn't use index keys, instead it partially unserialises each record in the results set
// - only unserialising the single required field for the aggregation (so its not too slow).
// For simplicity of code, API and to reduce reflection, the result returned is a float64,
// but Sum() will work on any number that is parsable from JSON as a float - so just convert to
// your required number type after the result is in.
// Sum() expects you to specify the path to the number of interest in your JSON using a string of field
// names representing the nested JSON path.  It's fairly intuitive,
// but see the docs for json parser (https://github.com/buger/jsonparser) for full details
func (q *Query) Sum(jsonPath []string) (float64, int, error) {
	var sum float64
	q.aggTarget = &sum
	q.slowSumPath = jsonPath
	n, err := q.execute()
	return sum, n, err
}

// Query Combination

// Or takes any number of queries and combines their results (as IDs) in a logical OR manner,
// returning one query, marked as executed, with union of IDs returned by the query.  The resulting query
// can be run, or combined further
func (db DB) Or(entities interface{}, queries ...*Query) *Query {
	return queryCombine(db, entities, union, queries...)
}

// Or takes any number of queries and combines their results (as IDs) in a logical AND manner,
// returning one query, marked as executed, with union of IDs returned by the query.  The resulting query
// can be run, or combined further
func (db DB) And(entities interface{}, queries ...*Query) *Query {
	return queryCombine(db, entities, intersection, queries...)
}
