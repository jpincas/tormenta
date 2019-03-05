package tormenta

import (
	"fmt"
	"strings"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/wsxiaoys/terminal/color"
)

func debugLogGet(target interface{}, start time.Time, noResults int, err error, ids ...gouuidv6.UUID) {
	if err != nil {
		msg := color.Sprintf("@Returned error: %s", err)
		fmt.Println(msg)
		return
	}

	entityName, _ := entityTypeAndValue(target)

	var idsStrings []string
	for _, id := range ids {
		idsStrings = append(idsStrings, id.String())
	}
	idsOutput := strings.Join(idsStrings, ",")

	msg := color.Sprintf(
		"@{!}GET@{|} @y[%s | %s]@{|} returned @c%v@{|} result(s) in @g%s",
		entityName,
		idsOutput,
		noResults,
		time.Since(start),
	)
	fmt.Println(msg)
}

func (q Query) debugLog(start time.Time, noResults int, err error) {
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

	msg := color.Sprintf("@{!}FIND@{|} @y[%s]@{|} returned @c%v@{|} result(s) in @g%s", q, noResults, time.Since(start))
	fmt.Println(msg)
}
