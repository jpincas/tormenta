package tormenta

import "sync"

type queryResult struct {
	ids idList
	err error
}

func queryCombine(combineFunc func(...idList) idList, queries ...*Query) *Query {
	firstQuery := queries[0]
	combinedQuery := &Query{
		db:            firstQuery.db,
		combinedQuery: true,
		target:        firstQuery.target,
		ctx:           firstQuery.ctx,
	}

	ch := make(chan queryResult)
	var wg sync.WaitGroup

	var queryIDs []idList
	var errorsList []error

	for _, query := range queries {
		// Regular, non-combined queries need to be run
		// through the id fether.  We fire those off in parallel
		if !query.combinedQuery {
			wg.Add(1)
			go func(thisQuery *Query) {
				err := thisQuery.queryIDs()
				ch <- queryResult{
					ids: thisQuery.ids,
					err: err,
				}
			}(query)
		} else {
			// Otherwise, if this is a nested combined query,
			// we can just add the list of ids as is
			queryIDs = append(queryIDs, query.ids)
		}
	}

	go func() {
		for queryResult := range ch {
			if queryResult.err != nil {
				errorsList = append(errorsList, queryResult.err)
			} else {
				queryIDs = append(queryIDs, queryResult.ids)
			}

			// Only signal to the wait group that a record has been fetched
			// at this point rather than the anonymous func above, otherwise
			// you tend to lose the last result
			wg.Done()
		}
	}()

	wg.Wait()

	combinedQuery.ids = combineFunc(queryIDs...)
	return combinedQuery
}
