package tormentarest

import (
	"fmt"
	"net/http"
	"strconv"

	tormenta "github.com/jpincas/tormenta/tormentadb"
)

const (
	limit   = "limit"
	reverse = "reverse"
)

func buildQuery(r *http.Request, q *tormenta.Query) error {
	// Reverse
	reverseString := r.URL.Query().Get(reverse)
	if reverseString == "true" {
		q.Reverse()
	} else if reverseString == "false" {
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

	return nil
}
