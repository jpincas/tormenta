package tormenta

import (
	"reflect"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

type Query struct {
	db       DB
	keyRoot  []byte
	value    reflect.Value
	entity   interface{}
	entities interface{}
	from, to gouuidv6.UUID
}

func (db DB) Query(entities interface{}) *Query {

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

	return &newQuery
}

func (q *Query) From(t time.Time) *Query {
	q.from = gouuidv6.NewFromTime(t)
	return q
}

func (q *Query) To(t time.Time) *Query {
	q.to = gouuidv6.NewFromTime(t)
	return q
}

func (q Query) rangePrefixes() (from, to []byte) {
	if !q.from.IsNil() {
		from = makePrefix(q.keyRoot, q.from.Bytes())
	} else {
		from = makePrefix(q.keyRoot, []byte{})
	}

	to = makePrefix(q.keyRoot, []byte{})

	return
}

func (q *Query) Run() (int, error) {
	// Work out what prefix to iterate over
	from, to := q.rangePrefixes()
	compareKey := compareKey(q.to, to)

	// Cast the 'entity' we have stored on the query to a 'tormentable'
	// so that it can be unmarshalled
	entity := q.entity.(tormentable)

	// Set up a slice to accumulate the results
	results := []tormentable{}

	err := q.db.KV.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(from); it.ValidForPrefix(to); it.Next() {

			key := it.Item().Key()
			if !q.to.IsNil() && !compare(compareKey, key) {
				break
			}

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
