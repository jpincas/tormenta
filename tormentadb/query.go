package tormentadb

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

	// List results placeholder
	rt reflect.Value

	// Counter
	counter int

	// Is this an aggregation query?
	isAggQuery bool
	aggTarget  interface{}
}

func (db DB) newQuery(target interface{}, first bool) *query {
	// Get the key root and cache the value
	keyRoot, value := entityTypeAndValue(target)

	// Create the base query
	q := &query{
		db:      db,
		keyRoot: keyRoot,
		value:   value,
		target:  target,
	}

	// Set up list results placeholder
	q.rt = reflect.Indirect(reflect.ValueOf(q.target))

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

func (q query) aggregate(item *badger.Item) {
	// TODO: super inefficient to do this every time
	switch q.aggTarget.(type) {
	case *int32:
		acc := *q.aggTarget.(*int32)
		extractIndexValue(item.Key(), q.aggTarget)
		*q.aggTarget.(*int32) = acc + *q.aggTarget.(*int32)
	case *float64:
		acc := *q.aggTarget.(*float64)
		extractIndexValue(item.Key(), q.aggTarget)
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

}

func (q *query) prepareQuery() {
	q.setRanges()
	q.resetQuery()
}

func (q *query) fetchIndexedRecord(item *badger.Item) error {
	key := extractID(item.Key())

	var entity Tormentable
	if q.first {
		entity = reflect.New(q.value.Type()).Interface().(Tormentable)
	} else {
		entity = reflect.New(q.value.Type().Elem()).Interface().(Tormentable)
	}

	// Get the record
	_, err := q.db.Get(entity, key)
	if err != nil {
		return err
	}

	// If this is a 'first' query - then just set the unmarshalled entity on the target
	// Otherwise, build up the results slice - we'll set on the target later!
	if q.first {
		reflect.Indirect(reflect.ValueOf(q.target)).Set(reflect.Indirect(reflect.ValueOf(entity)))
	} else {
		q.rt = reflect.Append(
			q.rt,
			reflect.Indirect(reflect.ValueOf(entity)),
		)
	}

	return nil
}

func (q *query) fetchRecord(item *badger.Item) error {
	// Set up the entity for unmarshalling
	var entity Tormentable
	if q.first {
		entity = reflect.New(q.value.Type()).Interface().(Tormentable)
	} else {
		entity = reflect.New(q.value.Type().Elem()).Interface().(Tormentable)
	}

	val, err := item.Value()
	if err != nil {
		return err
	}

	_, err = entity.UnmarshalMsg(val)
	if err != nil {
		return err
	}

	// If this is a 'first' query - then just set the unmarshalled entity on the target
	// Otherwise, build up the results slice - we'll set on the target later!
	if q.first {
		reflect.Indirect(reflect.ValueOf(q.target)).Set(reflect.Indirect(reflect.ValueOf(entity)))
	} else {
		q.rt = reflect.Append(
			q.rt,
			reflect.Indirect(reflect.ValueOf(entity)),
		)
	}

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
			item := it.Item()
			if !q.countOnly && !q.isAggQuery {
				if q.isIndexQuery {
					q.fetchIndexedRecord(item)
				} else {
					q.fetchRecord(item)
				}
			}

			if q.isAggQuery {
				q.aggregate(item)
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

	// If this was a first-only query
	if q.first {
		return q.counter, nil
	}

	// If this was a non-count-only and non-index query
	// we now need to set the results on the target
	if !q.isIndexQuery && !q.countOnly {
		reflect.Indirect(reflect.ValueOf(q.target)).Set(q.rt)
	}

	// Finally, return the numbrer of records found
	return q.counter, nil
}
