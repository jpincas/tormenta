package tormentarest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	tormenta "github.com/jpincas/tormenta/tormentadb"
)

const (
	limit   = "limit"
	reverse = "reverse"
	from    = "from"
	to      = "to"
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

	return nil
}
