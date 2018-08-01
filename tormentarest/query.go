package tormentarest

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	tormenta "github.com/jpincas/tormenta/tormentadb"
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
)

func buildQuery(r *http.Request, q *tormenta.Query) error {
	// Reverse
	reverseString := r.URL.Query().Get(reverse)
	if reverseString == "true" {
		q.Reverse()
	} else if reverseString == "false" || reverseString == "" {
		// don't actually do anything
		// but it's still a valid input
	} else {
		return fmt.Errorf(errBadReverseFormat, reverseString)
	}

	// Limit
	limitString := r.URL.Query().Get(limit)
	if limitString != "" {
		n, err := strconv.Atoi(limitString)
		if err != nil {
			return fmt.Errorf(errBadLimitFormat, limitString)
		}

		q.Limit(n)
	}

	// Offset
	offsetString := r.URL.Query().Get(offset)
	if offsetString != "" {
		n, err := strconv.Atoi(offsetString)
		if err != nil {
			return fmt.Errorf(errBadOffsetFormat, limitString)
		}

		q.Offset(n)
	}

	// From / To
	const dateFormat1 = "2006-01-02"
	fromString := r.URL.Query().Get(from)
	toString := r.URL.Query().Get(to)
	if fromString != "" {
		t, err := time.Parse(dateFormat1, fromString)
		if err != nil {
			return errors.New(errBadFromFormat)
		}
		q.From(t)
	}
	if toString != "" {
		t, err := time.Parse(dateFormat1, toString)
		if err != nil {
			return errors.New(errBadToFormat)
		}
		q.To(t)
	}

	// Index
	// Need to specify the index in a separate param
	// If an index param is found, we go into index query building mode
	indexString := r.URL.Query().Get(index)
	if indexString != "" {
		err := buildIndexQuery(r, q, indexString)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildIndexQuery(r *http.Request, q *tormenta.Query, key string) error {
	matchString := r.URL.Query().Get(match)
	startString := r.URL.Query().Get(start)
	endString := r.URL.Query().Get(end)

	// If 'index' param has been specified,
	// but no exact match or range has been given, return an error
	if matchString == "" && (startString == "" && endString == "") {
		return errors.New(errIndexWithNoParams)
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
			return errors.New(errRangeTypeMismatch)
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
