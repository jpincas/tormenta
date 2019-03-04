package tormenta

import (
	"fmt"
	"log"
	"time"
)

func (q Query) DebugLog(start time.Time, noResults int, err error) {
	// To log, we either need to be either in global debug mode,
	// or the debug flag for this query needs to be set to true
	if !q.db.Options.DebugMode && !q.debug {
		return
	}

	if err != nil {
		msg := fmt.Sprintf("Query returned error: %s", err)
		log.Println(msg)
		return
	}

	msg := fmt.Sprintf("[%s] returned %v results in %s", q, noResults, time.Since(start))
	log.Println(msg)
}
