package tormenta

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

type Query struct {
	// Connection to BadgerDB
	db DB

	// The entity type being searched
	keyRoot []byte

	// Target is the pointer passed into the Query where results will be set
	target interface{}

	single bool

	// Order by index name
	orderByIndexName []byte

	// Limit number of returned results
	limit int

	// Offet - start returning results N entities from the beginning
	// offsetCounter used to track the offset
	offset, offsetCounter int

	// Reverse fullStruct of searching and returned results
	reverse bool

	// From and To dates for the search
	from, to gouuidv6.UUID

	// Is this a count only search
	countOnly bool

	// A placeholders for errors to be passed down through the Query
	err error

	// Ranges and comparision key
	seekFrom, validTo, compareTo []byte

	sumIndexName []byte
	sumTarget    interface{}

	// Pass-through context
	ctx map[string]interface{}

	// Filter
	filters    []filter
	basicQuery *basicQuery

	// Logical ID combinator
	idsCombinator func(...idList) idList

	// Is already prepared?
	prepared bool
}

func (q Query) Compare(cq Query) bool {
	return fmt.Sprint(q) == fmt.Sprint(cq)
}

func fromUUID(t time.Time) gouuidv6.UUID {
	// Subtract 1 nanosecond form the specified time
	// Leads to an inclusive date search
	t = t.Add(-1 * time.Nanosecond)
	return gouuidv6.NewFromTime(t)
}

func toUUID(t time.Time) gouuidv6.UUID {
	return gouuidv6.NewFromTime(t)
}

func (db DB) newQuery(target interface{}) *Query {
	// Create the base Query
	q := &Query{
		db:      db,
		keyRoot: KeyRoot(target),
		target:  target,
	}

	// Start with blank context
	q.ctx = make(map[string]interface{})

	// Defualt to logical AND combination
	q.idsCombinator = intersection

	return q
}

func (q *Query) addFilter(f filter) {
	q.filters = append(q.filters, f)
}

func (q Query) shouldApplyLimitOffsetToFilter() bool {
	// We only pass the limit/offset to a filter if
	// there is only 1 filter AND there is no order by index
	return len(q.filters) == 1 && len(q.orderByIndexName) == 0
}

func (q Query) shouldApplyLimitOffsetToBasicQuery() bool {
	return len(q.orderByIndexName) == 0
}

func (q *Query) prepareQuery() {
	// Each filter also needs some of the top level information
	// e.g keyroot, date range, limit, offset etc,
	// so we copy that in now
	for i := range q.filters {
		q.filters[i].keyRoot = q.keyRoot
		q.filters[i].reverse = q.reverse
		q.filters[i].from = q.from
		q.filters[i].to = q.to

		if q.shouldApplyLimitOffsetToFilter() {
			q.filters[i].limit = q.limit
			q.filters[i].offset = q.offset
		}
	}

	// If there are no filters, then we prepare a 'basic query'
	if len(q.filters) == 0 {
		bq := &basicQuery{
			from:    q.from,
			to:      q.to,
			reverse: q.reverse,
			keyRoot: q.keyRoot,
		}

		if q.shouldApplyLimitOffsetToBasicQuery() {
			bq.limit = q.limit
			bq.offset = q.offset
		}

		q.basicQuery = bq
	}
}

func (q *Query) queryIDs(txn *badger.Txn) (idList, error) {
	if !q.prepared {
		q.prepareQuery()
	}

	var allResults []idList

	// If during the query planning and preparation,
	// something has gone wrong and an error has been set on the query,
	// we'll return right here and now
	if q.err != nil {
		return idList{}, q.err
	}

	if len(q.filters) > 0 {
		// FOR WHEN THERE ARE INDEX FILTERS
		// We process them serially at the moment, becuase Badger can only support 1 iterator
		// per transaction.  If that limitation is ever removed, we could do this in parallel
		for _, filter := range q.filters {
			thisFilterResults := filter.queryIDs(txn)
			allResults = append(allResults, thisFilterResults)
		}
	} else {
		// FOR WHEN THERE ARE NO INDEX FILTERS
		allResults = []idList{q.basicQuery.queryIDs(txn)}
	}

	// Combine the results from multiple filters,
	// or the single top level id list into one, final id list
	// according to the required AND/OR logic
	return q.idsCombinator(allResults...), nil
}

func (q *Query) execute() (int, error) {
	txn := q.db.KV.NewTransaction(false)
	defer txn.Discard()

	finalIDList, err := q.queryIDs(txn)
	if err != nil {
		return 0, err
	}

	// TODO: more conditions to restrict when this is necessary
	if len(q.orderByIndexName) > 0 {
		is := indexSearch{
			idsToSearchFor: finalIDList,
			reverse:        q.reverse,
			limit:          q.limit,
			keyRoot:        q.keyRoot,
			indexName:      q.orderByIndexName,
			offset:         q.offset,
		}

		// If we are doing a quicksum and the sum index is the same
		// as the order index, we can take advantage of this index
		// iteration to do the sum
		if len(q.sumIndexName) > 0 && q.sumTarget != nil {
			if string(q.sumIndexName) == string(q.orderByIndexName) {
				is.sumIndexName = q.sumIndexName
				is.sumTarget = q.sumTarget
			}
		}

		// This will order and apply limit/offset
		finalIDList = is.execute(txn)
	}

	// For count-only, there's nothing more to do
	if q.countOnly {
		return len(finalIDList), nil
	}

	// If a sumIndexName and a target have been specified,
	// then we will take that to mean that this is a quicksum execution
	// How we handle quicksum depends on wehther the sum index is different from the order index.
	// If the two are the same, then we have already worked out the quicksum in the index iteration above, and theres
	// no need to do it again
	if len(q.sumIndexName) > 0 && q.sumTarget != nil {
		if string(q.sumIndexName) != string(q.orderByIndexName) {
			is := indexSearch{
				idsToSearchFor: finalIDList,
				reverse:        q.reverse,
				limit:          q.limit,
				keyRoot:        q.keyRoot,
				indexName:      q.sumIndexName,
				offset:         q.offset,
				sumIndexName:   q.sumIndexName,
				sumTarget:      q.sumTarget,
			}

			is.execute(txn)
		}

		// Now, whether the quicksum was on the same index as order,
		// or any other index, we will have the result in the target, so we can return now
		return len(finalIDList), nil
	}

	// For 'First' type queries
	if q.single {
		// For 'first' queries, we should check that there is at least 1 record found
		// before trying to set it
		if len(finalIDList) == 0 {
			return 0, nil
		}

		// db.get ususally takes a 'Record', so we need to set a new one up
		// and then set the result of get to the target aftwards
		record := newRecord(q.target)
		id := finalIDList[0]
		if found, err := q.db.get(txn, record, q.ctx, id); err != nil {
			return 0, err
		} else if !found {
			return 0, fmt.Errorf("Could not retrieve record with id: %v", id)
		}

		setSingleResultOntoTarget(q.target, record)
		return 1, nil
	}

	// Otherwise we just get the records and return
	n, err := q.db.getIDsWithContext(txn, q.target, q.ctx, finalIDList...)
	if err != nil {
		return 0, err
	}

	return n, nil
}
