package tormenta

import (
	"reflect"

	"github.com/dgraph-io/badger"
)

type Query struct {
	db       DB
	keyRoot  []byte
	value    reflect.Value
	entity   interface{}
	entities interface{}
}

func (db DB) Query(entities interface{}) Query {

	newQuery := Query{}
	keyRoot, value := getKeyRoot(entities)
	newQuery.keyRoot = keyRoot
	newQuery.value = value
	newQuery.db = db
	newQuery.entities = entities

	// Get the underlying type of the elements of the 'entities' slice
	// cast to interface and save on the query.
	// This will be used for unmarhsalling
	newQuery.entity = reflect.New(value.Type().Elem()).Interface()

	return newQuery
}

func (q Query) Run() (int, error) {
	// Work out what prefix to iterate over
	prefix := makePrefix(q.keyRoot)

	// Cast the 'entity' we have stored on the query to a 'tormentable'
	// so that it can be unmarshalled
	entity := q.entity.(tormentable)

	// Set up a slice to accumulate the results
	results := []tormentable{}

	err := q.db.KV.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			val, err := it.Item().Value()
			if err != nil {
				return err
			}

			_, err = entity.UnmarshalMsg(val)
			if err != nil {
				return err
			}

			results = append(results, entity)
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	// Now we have a slice of 'tormentables'

	// Set up a slice for 'translating' the 'tormentables' into the target slice type
	rt := reflect.Indirect(reflect.ValueOf(q.entities))

	// Iterate through the result, using 'reflect.Append' to append to the above slice
	// the underlying type of the result
	for _, result := range results {
		rt = reflect.Append(
			rt,
			reflect.Indirect(reflect.ValueOf(result)),
		)
	}

	// Now set the accumulated, translated results on the original, passed in
	// 'entities'
	reflect.Indirect(reflect.ValueOf(q.entities)).Set(rt)

	return len(results), nil
}

// newQuery.entity = reflect.New(reflect.TypeOf(entities).Elem().Elem()).Interface()
