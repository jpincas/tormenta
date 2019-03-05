package tormenta

import (
	"fmt"
	"time"

	"github.com/wsxiaoys/terminal/color"
)

func (q Query) DebugLog(start time.Time, noResults int, err error) {
	// To log, we either need to be either in global debug mode,
	// or the debug flag for this query needs to be set to true
	if !q.db.Options.DebugMode && !q.debug {
		return
	}

	if err != nil {
		msg := color.Sprintf("@rQuery returned error: %s", err)
		fmt.Println(msg)
		return
	}

	msg := color.Sprintf("@{!}TORMENTA@{|} @y[%s]@{|} returned @c%v@{|} results in @g%s", q, noResults, time.Since(start))
	fmt.Println(msg)
}
