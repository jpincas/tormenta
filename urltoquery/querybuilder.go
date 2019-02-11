package urltoquery

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/jpincas/tormenta"
)

const (
	limit   = "limit"
	reverse = "reverse"
	from    = "from"
	to      = "to"
	offset  = "offset"

	// Index queries
	index = "index"
	match = "match"
	start = "start"
	end   = "end"

	ErrBadIDFormat                 = "Bad format for Tormenta ID - %s"
	ErrBadLimitFormat              = "%s is an invalid input for LIMIT. Expecting a number"
	ErrBadOffsetFormat             = "%s is an invalid input for OFFSET. Expecting a number"
	ErrBadReverseFormat            = "%s is an invalid input for REVERSE. Expecting true/false"
	ErrBadFromFormat               = "Invalid input for FROM. Expecting somthing like '2006-01-02'"
	ErrBadToFormat                 = "Invalid input for TO. Expecting somthing like '2006-01-02'"
	ErrFromIsAfterTo               = "FROM date is after TO date, making the date range impossible"
	ErrIndexWithNoParams           = "An index search has been specified, but no MATCH or START/END (for range) has been specified"
	ErrIndexMatchAndRangeSpecified = "An index search can be MATCH or RANGE, but not both"
	ErrParamsWithNoIndex           = "MATCH or START or END has been specified, but no index has been specified"
	ErrRangeTypeMismatch           = "For a range index search, START and END should be of the same type (bool, int, float, string)"
	ErrUnmarshall                  = "Error in format of data to save: %v"
)

// HasReverseBeenSet indicates whether the 'reverse' parameter has been specified
func HasReverseBeenSet(values url.Values) bool {
	reverseString := values.Get(reverse)
	if reverseString == "true" || reverseString == "false" {
		return true
	}

	return false
}

// BuildQuery builds up a passed-in query from url parameters.  It includes the option
// to ignore any limit param, which is useful if you want an overall count for a particular
// set of query conditions
func BuildQuery(q *tormenta.Query, values url.Values, ignoreLimitOffset bool) error {
	// Reverse
	reverseString := values.Get(reverse)
	if reverseString == "true" {
		q.Reverse()
	} else if reverseString == "false" || reverseString == "" {
		// don't actually do anything
		// but it's still a valid input
	} else {
		return fmt.Errorf(ErrBadReverseFormat, reverseString)
	}

	// Only apply limit and offset if required
	if !ignoreLimitOffset {
		// Limit
		limitString := values.Get(limit)
		if limitString != "" {
			n, err := strconv.Atoi(limitString)
			if err != nil {
				return fmt.Errorf(ErrBadLimitFormat, limitString)
			}

			q.Limit(n)
		}

		// Offset
		offsetString := values.Get(offset)
		if offsetString != "" {
			n, err := strconv.Atoi(offsetString)
			if err != nil {
				return fmt.Errorf(ErrBadOffsetFormat, offsetString)
			}

			q.Offset(n)
		}
	}

	// From / To
	const dateFormat1 = "2006-01-02"
	fromString := values.Get(from)
	toString := values.Get(to)

	var toValue, fromValue time.Time
	var err error

	if fromString != "" {
		fromValue, err = time.Parse(dateFormat1, fromString)
		if err != nil {
			return errors.New(ErrBadFromFormat)
		}
		q.From(fromValue)
	}

	if toString != "" {
		toValue, err = time.Parse(dateFormat1, toString)
		if err != nil {
			return errors.New(ErrBadToFormat)
		}
		q.To(toValue)
	}

	// If both from and to where specified, make sure to is later
	if fromString != "" && toString != "" && fromValue.After(toValue) {
		return errors.New(ErrFromIsAfterTo)
	}

	// Index
	// Need to specify the index in a separate param
	// If an index param is found, we go into index query building mode
	indexString := values.Get(index)
	if indexString != "" {
		err := buildIndexQuery(q, values, indexString)
		if err != nil {
			return err
		}
	} else {
		// If no index param is found, then we should check for match, start, end and error if they are found
		matchString := values.Get(match)
		startString := values.Get(start)
		endString := values.Get(end)
		if matchString != "" || startString != "" || endString != "" {
			return errors.New(ErrParamsWithNoIndex)
		}

	}

	return nil
}

func buildIndexQuery(q *tormenta.Query, values url.Values, key string) error {
	matchString := values.Get(match)
	startString := values.Get(start)
	endString := values.Get(end)

	// if no exact match or range has been given, return an error
	if matchString == "" && (startString == "" && endString == "") {
		return errors.New(ErrIndexWithNoParams)
	}

	// likewise, if match AND a range (either start or end) have been specified
	if matchString != "" && (startString != "" || endString != "") {
		return errors.New(ErrIndexMatchAndRangeSpecified)
	}

	// Match
	// We'll process this first, so if a 'match' clause is specified
	// AND range is also specified, match will take precedence
	if matchString != "" {
		q.Match(key, stringToInterface(matchString))
		return nil
	}

	// Range
	// If both START and END are specified,
	// they should be of the same type
	if startString != "" && endString != "" {
		start := stringToInterface(startString)
		end := stringToInterface(endString)
		if reflect.TypeOf(start) !=
			reflect.TypeOf(end) {
			return errors.New(ErrRangeTypeMismatch)
		}

		q.Range(key, start, end)
		return nil
	}

	// START only
	if startString != "" {
		q.Range(key, stringToInterface(startString), nil)
		return nil
	}

	// END only
	if endString != "" {
		q.Range(key, nil, stringToInterface(endString))
		return nil
	}

	return nil
}

func stringToInterface(s string) interface{} {
	// Int
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	}

	// Float
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}

	// Bool
	// Bool last, otherwise 0/1 get wrongly interpreted
	b, err := strconv.ParseBool(s)
	if err == nil {
		return b
	}

	// Default to string
	return s
}
