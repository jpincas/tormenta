package urltoquery

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jpincas/tormenta"
)

type (
	queryString string
	queryValues map[string]string
)

func (qs queryString) parse() (queryValues, error) {
	queryValues := queryValues{}

	components := strings.Split(string(qs), qsSeparator)
	for _, component := range components {
		queryKV := strings.Split(component, qvSeparator)
		if len(queryKV) != 2 {
			return queryValues, errors.New(ErrBadFormatQueryValue)
		}

		queryValues[queryKV[0]] = queryKV[1]
	}

	return queryValues, nil
}

func (qv queryValues) get(key string) string {
	return qv[key]
}

func (qs queryString) toQuery(db *tormenta.DB, target interface{}, ignoreLimitOffset bool) (*tormenta.Query, error) {
	// Set up the base query
	q := db.Find(target)

	values, err := qs.parse()
	if err != nil {
		return nil, err
	}

	// Reverse
	reverseString := values.get(reverse)
	if reverseString == "true" {
		q.Reverse()
	} else if reverseString == "false" || reverseString == "" {
		// don't actually do anything
		// but it's still a valid input
	} else {
		return nil, fmt.Errorf(ErrBadReverseFormat, reverseString)
	}

	// Only apply limit and offset if required
	if !ignoreLimitOffset {
		// Limit
		limitString := values.get(limit)

		if limitString != "" {
			n, err := strconv.Atoi(limitString)
			if err != nil {
				return nil, fmt.Errorf(ErrBadLimitFormat, limitString)
			}

			q.Limit(n)
		}

		// Offset
		offsetString := values.get(offset)
		if offsetString != "" {
			n, err := strconv.Atoi(offsetString)
			if err != nil {
				return nil, fmt.Errorf(ErrBadOffsetFormat, offsetString)
			}

			q.Offset(n)
		}
	}

	// From / To
	const dateFormat1 = "2006-01-02"
	fromString := values.get(from)
	toString := values.get(to)

	var toValue, fromValue time.Time

	if fromString != "" {
		fromValue, err = time.Parse(dateFormat1, fromString)
		if err != nil {
			return nil, errors.New(ErrBadFromFormat)
		}
		q.From(fromValue)
	}

	if toString != "" {
		toValue, err = time.Parse(dateFormat1, toString)
		if err != nil {
			return nil, errors.New(ErrBadToFormat)
		}
		q.To(toValue)
	}

	// If both from and to where specified, make sure to is later
	if fromString != "" && toString != "" && fromValue.After(toValue) {
		return nil, errors.New(ErrFromIsAfterTo)
	}

	// Index
	// Need to specify the index in a separate param
	// If an index param is found, we go into index query building mode
	indexString := values.get(index)
	if indexString != "" {
		err := buildIndexQuery(q, values, indexString)
		if err != nil {
			return nil, err
		}
	} else {
		// If no index param is found, then we should check for match, start, end and error if they are found
		matchString := values.get(match)
		startString := values.get(start)
		endString := values.get(end)
		if matchString != "" || startString != "" || endString != "" {
			return nil, errors.New(ErrParamsWithNoIndex)
		}

	}

	return q, nil
}

func buildIndexQuery(q *tormenta.Query, values queryValues, key string) error {
	matchString := values.get(match)
	startsWithString := values.get(startsWith)
	startString := values.get(start)
	endString := values.get(end)

	// if no exact match or range or starsWith has been given, return an error
	if matchString == "" && startsWithString == "" && (startString == "" && endString == "") {
		return errors.New(ErrIndexWithNoParams)
	}

	// If more than one of MATCH, RANGE and STARTSWITH have been specified
	if matchString != "" && (startString != "" || endString != "") ||
		matchString != "" && startsWithString != "" ||
		startsWithString != "" && (startString != "" || endString != "") {
		return errors.New(ErrTooManyIndexOperatorsSpecified)
	}

	if matchString != "" {
		q.Match(key, stringToInterface(matchString))
		return nil
	}

	if startsWithString != "" {
		q.StartsWith(key, startsWithString)
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
