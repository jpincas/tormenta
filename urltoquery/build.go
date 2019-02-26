package urltoquery

import (
	"net/url"

	"github.com/jpincas/tormenta"
)

const (
	// Top level labels that can be parsed out of the URL string
	query = "query"
	or    = "or"

	// Symbols used in specifying query key:value pairs
	// INSIDE a url param e.g. query=myKey:myValue,anotherKey:anotherValue
	qvSeparator = ":"
	qsSeparator = ","

	// Keywords for query parameters
	offset  = "offset"
	limit   = "limit"
	reverse = "reverse"
	from    = "from"
	to      = "to"

	// Keywords for index query components
	index      = "index"
	match      = "match"
	startsWith = "starts"
	start      = "start"
	end        = "end"

	// Error messages
	ErrBadFormatQueryValue            = "Bad format for query value"
	ErrBadIDFormat                    = "Bad format for Tormenta ID - %s"
	ErrBadLimitFormat                 = "%s is an invalid input for LIMIT. Expecting a number"
	ErrBadOffsetFormat                = "%s is an invalid input for OFFSET. Expecting a number"
	ErrBadReverseFormat               = "%s is an invalid input for REVERSE. Expecting true/false"
	ErrBadFromFormat                  = "Invalid input for FROM. Expecting somthing like '2006-01-02'"
	ErrBadToFormat                    = "Invalid input for TO. Expecting somthing like '2006-01-02'"
	ErrFromIsAfterTo                  = "FROM date is after TO date, making the date range impossible"
	ErrIndexWithNoParams              = "An index search has been specified, but no MATCH or START/END (for range) has been specified"
	ErrTooManyIndexOperatorsSpecified = "An index search can be MATCH, RANGE or STARTSWITH, but not multiple matching operators"
	ErrParamsWithNoIndex              = "MATCH or START or END has been specified, but no index has been specified"
	ErrRangeTypeMismatch              = "For a range index search, START and END should be of the same type (bool, int, float, string)"
	ErrUnmarshall                     = "Error in format of data to save: %v"
)

// HasReverseBeenSet indicates whether the 'reverse' parameter has been specified
func HasReverseBeenSet(values url.Values) bool {
	reverseString := values.Get(reverse)
	if reverseString == "true" || reverseString == "false" {
		return true
	}

	return false
}

func BuildQuery(db *tormenta.DB, target interface{}, values url.Values, ignoreLimitOffset bool) (*tormenta.Query, error) {
	queries, err := buildQueries(db, target, values, ignoreLimitOffset)
	if err != nil {
		return nil, err
	}

	return combineQueries(db, target, values, queries)
}

func combineQueries(db *tormenta.DB, target interface{}, values url.Values, queries []*tormenta.Query) (*tormenta.Query, error) {
	// Default to and AND combinator
	orString := values.Get(or)
	if orString == "true" {
		return db.Or(target, queries...), nil
	}

	return db.And(target, queries...), nil
}

// BuildQuery builds up a passed-in query from url parameters.  It includes the option
// to ignore any limit param, which is useful if you want an overall count for a particular
// set of query conditions
func buildQueries(db *tormenta.DB, target interface{}, values url.Values, ignoreLimitOffset bool) ([]*tormenta.Query, error) {
	// Placeholder for the queries that will be combined at the end
	var queries []*tormenta.Query

	// Get each query from the url, parse and add to the query list
	queryStrings := values["query"]
	for _, qs := range queryStrings {
		query, err := queryString(qs).toQuery(db, target, ignoreLimitOffset)
		if err != nil {
			return nil, err
		}

		queries = append(queries, query)
	}

	return queries, nil
}
