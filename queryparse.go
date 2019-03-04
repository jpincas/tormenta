package tormenta

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// Symbols used in specifying query key:value pairs
	// INSIDE a url param e.g. query=myKey:myValue,anotherKey:anotherValue
	whereValueSeparator  = ":"
	whereClauseSeparator = ","

	queryStringWhere      = "where"
	queryStringOr         = "or"
	queryStringOrderBy    = "order"
	queryStringOffset     = "offset"
	queryStringLimit      = "limit"
	queryStringReverse    = "reverse"
	queryStringFrom       = "from"
	queryStringTo         = "to"
	queryStringMatch      = "match"
	queryStringStartsWith = "startswith"
	queryStringStart      = "start"
	queryStringEnd        = "end"
	queryStringIndex      = "index"

	// Error messages
	ErrBadFormatQueryValue            = "Bad format for query value"
	ErrBadIDFormat                    = "Bad format for Tormenta ID - %s"
	ErrBadLimitFormat                 = "%s is an invalid input for LIMIT. Expecting a number"
	ErrBadOffsetFormat                = "%s is an invalid input for OFFSET. Expecting a number"
	ErrBadReverseFormat               = "%s is an invalid input for REVERSE. Expecting true/false"
	ErrBadOrFormat                    = "%s is an invalid input for OR. Expecting true/false"
	ErrBadFromFormat                  = "Invalid input for FROM. Expecting somthing like '2006-01-02'"
	ErrBadToFormat                    = "Invalid input for TO. Expecting somthing like '2006-01-02'"
	ErrFromIsAfterTo                  = "FROM date is after TO date, making the date range impossible"
	ErrIndexWithNoParams              = "An index search has been specified, but index search operator has been specified"
	ErrTooManyIndexOperatorsSpecified = "An index search can be MATCH, RANGE or STARTSWITH, but not multiple matching operators"
	ErrWhereClauseNoIndex             = "A WHERE clause requires an index to be specified"
	ErrRangeTypeMismatch              = "For a range index search, START and END should be of the same type (bool, int, float, string)"
	ErrUnmarshall                     = "Error in format of data to save: %v"
)

func (q *Query) Parse(ignoreLimitOffset bool, s string) error {
	// Parse the query string for values
	values, err := url.ParseQuery(s)
	if err != nil {
		return err
	}

	// Reverse
	reverseString := values.Get(queryStringReverse)
	if reverseString == "true" {
		q.Reverse()
	} else if reverseString == "false" || reverseString == "" {
		q.UnReverse()
	} else {
		return fmt.Errorf(ErrBadReverseFormat, reverseString)
	}

	// Order by
	orderByString := values.Get(queryStringOrderBy)
	if orderByString != "" {
		q.OrderBy(orderByString)
	}

	// Only apply limit and offset if required
	if !ignoreLimitOffset {
		limitString := values.Get(queryStringLimit)

		if limitString != "" {
			n, err := strconv.Atoi(limitString)
			if err != nil {
				return fmt.Errorf(ErrBadLimitFormat, limitString)
			}

			q.Limit(n)
		}

		// Offset
		offsetString := values.Get(queryStringOffset)
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
	fromString := values.Get(queryStringFrom)
	toString := values.Get(queryStringTo)

	var toValue, fromValue time.Time

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

	// Process each where clause individually
	whereClauseStrings := values["where"]
	for _, w := range whereClauseStrings {
		if err := whereClauseString(w).addToQuery(q); err != nil {
			return err
		}
	}

	// And -> Or
	orString := values.Get(queryStringOr)
	if orString == "true" {
		q.Or()
	} else if orString == "false" {
		q.And() // this is the default anyway
	} else if orString == "" {
		// Nothing to do here
	} else {
		return fmt.Errorf(ErrBadOrFormat, orString)
	}

	return nil
}

type (
	whereClauseString string
	whereClauseValues map[string]string
)

func (wcs whereClauseString) parse() (whereClauseValues, error) {
	whereClauseValues := whereClauseValues{}

	components := strings.Split(string(wcs), whereClauseSeparator)
	for _, component := range components {
		whereKV := strings.Split(component, whereValueSeparator)
		if len(whereKV) != 2 {
			return whereClauseValues, errors.New(ErrBadFormatQueryValue)
		}

		whereClauseValues[whereKV[0]] = whereKV[1]
	}

	return whereClauseValues, nil
}

func (wcs whereClauseString) addToQuery(q *Query) error {
	values, err := wcs.parse()
	if err != nil {
		return err
	}

	indexString := values.get(queryStringIndex)
	if indexString == "" {
		return errors.New(ErrWhereClauseNoIndex)
	}

	return values.addToQuery(q, indexString)
}

func (values whereClauseValues) get(key string) string {
	return values[key]
}

func (values whereClauseValues) addToQuery(q *Query, key string) error {
	matchString := values.get(queryStringMatch)
	startsWithString := values.get(queryStringStartsWith)
	startString := values.get(queryStringStart)
	endString := values.get(queryStringEnd)

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

// String methods for queries and filters

type queryComponent struct {
	key   string
	value interface{}
}

func (q Query) String() string {

	isOr := isOr(q.idsCombinator)

	components := []queryComponent{
		{queryStringLimit, q.limit},
		{queryStringOffset, q.offset},
		{queryStringReverse, q.reverse},
		{queryStringOr, isOr},
		{queryStringFrom, q.from.Time()},
		{queryStringTo, q.to.Time()},
		{queryStringOrderBy, string(q.orderByIndexName)},
	}

	var componentStrings []string
	for _, component := range components {
		componentStrings = append(componentStrings, fmt.Sprintf("%s=%v", component.key, component.value))
	}

	for _, filter := range q.filters {
		componentStrings = append(componentStrings, fmt.Sprintf("%s=%s", queryStringWhere, filter.String()))
	}

	builtQuery := strings.Join(componentStrings, "&")

	return fmt.Sprintf("/%s?%s", string(q.keyRoot), builtQuery)
}

func (f filter) String() string {
	components := []queryComponent{
		{queryStringStart, f.start},
		{queryStringEnd, f.end},
		{queryStringIndex, string(f.indexName)},
		{queryStringStartsWith, f.isStartsWithQuery},
	}

	var componentStrings []string
	for _, component := range components {
		componentStrings = append(componentStrings, fmt.Sprintf("%s%s%v", component.key, whereValueSeparator, component.value))
	}

	return strings.Join(componentStrings, whereClauseSeparator)
}
