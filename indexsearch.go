package tormenta

import (
	"reflect"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

type indexQuery struct {
	db         DB
	keyRoot    []byte
	value      reflect.Value
	holder     interface{}
	target     interface{}
	limit      int
	reverse    bool
	first      bool
	start, end interface{}
	from, to   gouuidv6.UUID
	indexName  []byte
	countOnly  bool
}

func (db DB) newIndexQuery(target interface{}, first bool, indexName string) *indexQuery {
	// Get the key root and cache the value
	keyRoot, value := entityTypeAndValue(target)

	// Set the 'holder' on the query
	// This is basically the type of either the single entity passed in (for first only)
	// or the underlying element type for the slice (regular query)
	// During query execution this will be cast to a 'tormentable'
	// which will be used as a 'holder' for each unmarshalling operation
	var holder interface{}
	if first {
		holder = reflect.New(value.Type()).Interface()
	} else {
		holder = reflect.New(value.Type().Elem()).Interface()
	}

	q := &indexQuery{
		db:        db,
		keyRoot:   keyRoot,
		value:     value,
		holder:    holder,
		target:    target,
		indexName: []byte(indexName),
	}

	if first {
		q.limit = 1
		q.first = true
	}

	return q
}

func (db DB) FindIndex(entities interface{}, indexName string) *indexQuery {
	return db.newIndexQuery(entities, false, indexName)
}

func (q *indexQuery) Start(t interface{}) *indexQuery {
	q.start = t
	return q
}

func (q *indexQuery) End(t interface{}) *indexQuery {
	q.end = t
	return q
}

// From adds a lower boundary to the date range of the query
func (q *indexQuery) From(t time.Time) *indexQuery {
	q.from = gouuidv6.NewFromTime(t)
	return q
}

// To adds an upper bound to the date range of the query
func (q *indexQuery) To(t time.Time) *indexQuery {
	q.to = gouuidv6.NewFromTime(t)
	return q
}

// Match is shorthand for specifying Start() and End()
func (q *indexQuery) Match(t interface{}) *indexQuery {
	q.start = t
	q.end = t
	return q
}

func (q *indexQuery) Run() (int, error) {
	return q.execute()
}

func (q indexQuery) getIteratorOptions() badger.IteratorOptions {
	options := badger.DefaultIteratorOptions
	options.Reverse = q.reverse
	// Never prefetch values for index search
	options.PrefetchValues = false
	return options
}

func (q indexQuery) isExactIndexMatchSearch() bool {
	return q.start == q.end
}

func (q *indexQuery) execute() (int, error) {
	// If the end is the same as the start,
	// we can 'hardcode' the index value into the 'validTo' prefix
	// and therefore don't need to do any key comparision
	var seekFrom, validTo, compareTo []byte

	isExactMatchSearch := q.isExactIndexMatchSearch()

	if isExactMatchSearch {
		seekFrom = newIndexMatchKey(q.keyRoot, q.indexName, q.start, q.from).bytes()
		validTo = newIndexMatchKey(q.keyRoot, q.indexName, q.end).bytes()
		compareTo = newIndexMatchKey(q.keyRoot, q.indexName, q.end, q.to).bytes()
	} else {
		seekFrom = newIndexKey(q.keyRoot, q.indexName, q.start).bytes()
		validTo = newIndexKey(q.keyRoot, q.indexName, nil).bytes()
		compareTo = newIndexKey(q.keyRoot, q.indexName, q.end).bytes()
	}

	entity := q.holder.(Tormentable)
	results := []Tormentable{}
	counter := 0

	err := q.db.KV.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(q.getIteratorOptions())
		defer it.Close()

		for it.Seek(seekFrom); it.ValidForPrefix(validTo); it.Next() {

			// If a limit clause has been specified AND met, stop here
			if q.limit > 0 && counter >= q.limit {
				return nil
			}

			// Get the current key
			key := it.Item().Key()

			// Normally the final ID will require stripping for index searched
			// However, in the case of an exact match search, with a To clause, it should be left in
			stripID := true
			if isExactMatchSearch && !q.to.IsNil() {
				stripID = false
			}

			if q.end != nil && !compareKeyBytes(compareTo, key, q.reverse, stripID) {
				return nil
			}

			// Extract the ID from the key, and use to get the record by ID
			_, err := q.db.GetByID(entity, extractID(key))
			if err != nil {
				return err
			}

			// Append the retrieved record to the list of results
			results = append(results, entity)

			// Increment the counter
			counter++
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return len(results), nil
}
