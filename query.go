package tormenta

import (
	"reflect"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

type query struct {
	// Connection to BadgerDB
	db DB

	// The entity type being searched
	keyRoot []byte

	// The Go 'value' for the entity type being searched
	value reflect.Value

	// The 'holder' is basically the type of either the single entity passed in (for first only)
	// or the underlying element type for the slice (regular query).
	// During query execution this will be cast to a 'tormentable'
	// which will be used as a 'holder' for each unmarshalling operation
	holder interface{}

	// Target is the pointer passed into the query where results will be set
	target interface{}

	// Limit number of returned results
	limit int

	// Offet - start returning results N entities from the beginning
	// offsetCounter used to track the offset
	offset, offsetCounter int

	// Reverse order of searching and returned results
	reverse bool

	// Is this a 'first only' search
	first bool

	// The start and end points of the index range search
	start, end interface{}

	// From and To dates for the search
	from, to gouuidv6.UUID

	// If this is an index search, this is the name of the index
	indexName []byte

	// Is this an index query
	isIndexQuery bool

	// Is this a count only search
	countOnly bool

	// A placeholders for errors to be passed down through the query
	err error

	// Ranges and comparision key
	seekFrom, validTo, compareTo []byte

	// Entity placeholder
	entity Tormentable

	// Results holder
	results []Tormentable

	// Counter
	counter int

	// Is this an aggregation query?
	isAggQuery bool
	aggTarget  interface{}
}

func (db DB) newQuery(target interface{}, first bool) *query {
	// Get the key root and cache the value
	keyRoot, value := entityTypeAndValue(target)

	// Set the 'holder' on the query
	var holder interface{}
	if first {
		holder = reflect.New(value.Type()).Interface()
	} else {
		holder = reflect.New(value.Type().Elem()).Interface()
	}

	// Create the base query
	q := &query{
		db:      db,
		keyRoot: keyRoot,
		value:   value,
		holder:  holder,
		target:  target,
		entity:  holder.(Tormentable),
	}

	// If this is a 'first only' query
	if first {
		q.limit = 1
		q.first = true
	}

	return q
}

func (q query) getIteratorOptions(getValues bool) badger.IteratorOptions {
	options := badger.DefaultIteratorOptions
	options.Reverse = q.reverse
	options.PrefetchValues = getValues
	return options
}

func (q query) isExactIndexMatchSearch() bool {
	return q.start == q.end
}

func (q query) shouldGetValues() bool {
	// For index queries or count only queries, don't get values
	if q.isIndexQuery || q.countOnly {
		return false
	}

	return true
}

func (q query) shouldStripKeyID() bool {
	// Regular queries never need to have ID stripped
	if !q.isIndexQuery {
		return false
	}

	// Index queries which are exact match AND have a 'to' clause
	// also never need to have ID stripped
	if q.isExactIndexMatchSearch() && !q.to.IsNil() {
		return false
	}

	return true
}

func (q query) isEndOfRange(it *badger.Iterator) bool {
	key := it.Item().Key()

	if q.isIndexQuery {
		return q.end != nil && compareKeyBytes(q.compareTo, key, q.reverse, q.shouldStripKeyID())
	}

	return !q.to.IsNil() && compareKeyBytes(q.compareTo, key, q.reverse, q.shouldStripKeyID())
}

func (q query) isLimitMet() bool {
	return q.limit > 0 && q.counter >= q.limit
}

func (q query) endIteration(it *badger.Iterator) bool {
	if it.ValidForPrefix(q.validTo) {
		if q.isLimitMet() || q.isEndOfRange(it) {
			return false
		}

		return true
	}

	return false
}

func (q query) aggregate(it *badger.Iterator) {
	// TODO: super inefficient to do this every time
	switch q.aggTarget.(type) {
	case *int32:
		acc := *q.aggTarget.(*int32)
		extractIndexValue(it.Item().Key(), q.aggTarget)
		*q.aggTarget.(*int32) = acc + *q.aggTarget.(*int32)
	case *float64:
		acc := *q.aggTarget.(*float64)
		extractIndexValue(it.Item().Key(), q.aggTarget)
		*q.aggTarget.(*float64) = acc + *q.aggTarget.(*float64)
	}
}

func (q *query) setRanges() {
	var seekFrom, validTo, compareTo []byte

	if q.isIndexQuery && q.isExactIndexMatchSearch() {
		// For index searches with exact match
		seekFrom = newIndexMatchKey(q.keyRoot, q.indexName, q.start, q.from).bytes()
		validTo = newIndexMatchKey(q.keyRoot, q.indexName, q.end).bytes()
		compareTo = newIndexMatchKey(q.keyRoot, q.indexName, q.end, q.to).bytes()
	} else if q.isIndexQuery {
		// For regular index searches
		seekFrom = newIndexKey(q.keyRoot, q.indexName, q.start).bytes()
		validTo = newIndexKey(q.keyRoot, q.indexName, nil).bytes()
		compareTo = newIndexKey(q.keyRoot, q.indexName, q.end).bytes()
	} else {
		seekFrom = newContentKey(q.keyRoot, q.from).bytes()
		validTo = newContentKey(q.keyRoot).bytes()
		compareTo = newContentKey(q.keyRoot, q.to).bytes()
	}

	q.seekFrom = seekFrom
	q.validTo = validTo
	q.compareTo = compareTo
}

func (q *query) resetQuery() {
	// Counter should always be reset before executing a query.
	// Just in case a query is built then executed twice.
	q.counter = 0
	q.offsetCounter = q.offset
	q.results = []Tormentable{}
}

func (q *query) prepareQuery() {
	q.setRanges()
	q.resetQuery()
}

func (q *query) fetchIndexedRecord(it *badger.Iterator) error {
	key := extractID(it.Item().Key())

	// Get the record
	_, err := q.db.Get(q.entity, key)
	if err != nil {
		return err
	}

	// Append the retrieved record to the list of results
	q.results = append(q.results, q.entity)

	return nil
}

func (q *query) setFirst() {
	// The unmarhsalled entity will currently be on the query's 'entity' placeholder
	// So all we need to do now is set it on the target
	if q.isIndexQuery {
		q.target = q.entity
	}

	// For a regular query, the unmarhsalled entity will be on entity
	reflect.Indirect(reflect.ValueOf(q.target)).Set(reflect.Indirect(reflect.ValueOf(q.entity)))
}

func (q *query) fetchRecord(it *badger.Iterator) error {
	val, err := it.Item().Value()
	if err != nil {
		return err
	}

	_, err = q.entity.UnmarshalMsg(val)
	if err != nil {
		return err
	}

	q.results = append(q.results, q.entity)
	return nil
}

func (q *query) execute() (int, error) {
	// Do the work of calculating and setting initial values for the query
	q.prepareQuery()

	// Iterate through records according to calcuted range limits
	err := q.db.KV.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(q.getIteratorOptions(q.shouldGetValues()))
		defer it.Close()

		// Start iteration
		for it.Seek(q.seekFrom); q.endIteration(it); it.Next() {
			// If this is a 'range index' type query
			// that ALSO has a date range, the procedure is a little more complicated
			// compared to an exact index match.
			// Since the start/end points of the iteration focus on the index, e.g. E-J (alphabetical index)
			// we need to manually check all the keys and reject those that don't fit the date range
			if q.isIndexQuery && !q.isExactIndexMatchSearch() {
				key := extractID(it.Item().Key())
				if keyIsOutsideDateRange(key, q.from, q.to) {
					continue
				}
			}

			// Skip the first N entities according to the specified offset
			if q.offsetCounter > 0 {
				q.offsetCounter--
				continue
			}

			q.counter++

			// For non-count-only queries, we'll actually get the record
			// How this is done depends on whether this is an index-based search or not
			if !q.countOnly && !q.isAggQuery {
				if q.isIndexQuery {
					q.fetchIndexedRecord(it)
				} else {
					q.fetchRecord(it)
				}
			}

			if q.isAggQuery {
				q.aggregate(it)
			}

			// If this is a first-only search, break out of the iteration now
			// The counter has been incremented, so will read 1
			if q.first {
				return nil
			}
		}

		return nil
	})

	// If there was an error from the DB transaction, return 0 now
	if err != nil {
		return 0, err
	}

	// For count-only queries, there's nothing more to do
	if q.countOnly {
		return q.counter, nil
	}

	// If this was a first-only query, set the entity to the target
	// and return the counter value.
	// If no entity was found, the set will basically do nothing
	// and the counter will read 0
	if q.first {
		q.setFirst()
		return q.counter, nil
	}

	// If this was a non-count-only and non-index query
	// we now need to set the results on the target
	if !q.isIndexQuery && !q.countOnly {
		// Now we have a slice of 'Tormentables'
		// Set up a slice for 'translating' the 'Tormentables' into the target slice type
		rt := reflect.Indirect(reflect.ValueOf(q.target))

		// Iterate through the result, using 'reflect.Append' to append to the above slice
		// the underlying type of the result
		for _, result := range q.results {
			rt = reflect.Append(
				rt,
				reflect.Indirect(reflect.ValueOf(result)),
			)
		}

		// Now set the accumulated, translated results on the original, passed in
		// 'entities'
		reflect.Indirect(reflect.ValueOf(q.target)).Set(rt)
	}

	// Finally, return the numbrer of records found
	return len(q.results), nil
}
