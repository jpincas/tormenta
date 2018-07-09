package tormenta

import (
	"reflect"

	"github.com/dgraph-io/badger"
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

func (q *indexQuery) execute() (int, error) {
	seekFrom := newIndexKey(q.keyRoot, q.indexName, q.start).bytes()
	validTo := newIndexKey(q.keyRoot, q.indexName, nil).bytes()
	compareTo := newIndexKey(q.keyRoot, q.indexName, q.end).bytes()

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

			// Compare the current key to the calculated 'comparison' key
			// Specify 'true' so that the ID is stripped for comparison
			if q.end != nil && !compareKeyBytes(compareTo, key, q.reverse, true) {
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
