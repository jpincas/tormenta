package tormenta

import (
	"reflect"

	"github.com/dgraph-io/badger"
)

type indexQuery struct {
	db        DB
	keyRoot   []byte
	value     reflect.Value
	holder    interface{}
	target    interface{}
	limit     int
	reverse   bool
	first     bool
	from, to  interface{}
	indexName []byte
	countOnly bool
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

func (q *indexQuery) From(t interface{}) *indexQuery {
	q.from = t
	return q
}

func (q *indexQuery) To(t interface{}) *indexQuery {
	q.to = t
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
	seekFrom := newIndexKey(q.keyRoot, q.indexName, q.from).bytes()
	validTo := newIndexKey(q.keyRoot, q.indexName, nil).bytes()
	compareTo := newIndexKey(q.keyRoot, q.indexName, q.to).bytes()

	// entity := q.holder.(Tormentable)
	keys := [][]byte{}
	// results := []Tormentable{}
	counter := 0

	err := q.db.KV.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(q.getIteratorOptions())
		defer it.Close()

		for it.Seek(seekFrom); it.ValidForPrefix(validTo); it.Next() {

			if q.limit > 0 && counter >= q.limit {
				return nil
			}

			key := it.Item().Key()

			if q.to != nil && !compareKeyBytes(compareTo, key, q.reverse) {
				return nil
			}

			keys = append(keys, key)
			counter++
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return len(keys), nil
}
