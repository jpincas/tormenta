package tormenta

import "sync"

type queryResult struct {
	ids idList
	err error
}

func queryCombine(combineFunc func(...idList) idList, queries ...*Query) *Query {
	firstQuery := queries[0]
	combinedQuery := &Query{
		db:              firstQuery.db,
		alreadyExecuted: true,
		target:          firstQuery.target,
		ctx:             firstQuery.ctx,
	}

	ch := make(chan queryResult)
	var wg sync.WaitGroup

	for _, query := range queries {
		wg.Add(1)

		if !query.alreadyExecuted {
			go func(thisQuery *Query) {
				err := thisQuery.queryIDs()
				ch <- queryResult{
					ids: thisQuery.ids,
					err: err,
				}
			}(query)
		}
	}

	var queryIDs []idList
	var errorsList []error

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
