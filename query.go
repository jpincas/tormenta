package tormenta

import (
	"reflect"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

type query struct {
	db       DB
	keyRoot  []byte
	value    reflect.Value
	entity   interface{}
	entities interface{}
	from, to gouuidv6.UUID
	limit    int
	reverse  bool
}

func (db DB) Find(entities interface{}) *query {

	newquery := query{}
	keyRoot, value := getKeyRoot(entities)
	newquery.keyRoot = keyRoot
	newquery.value = value
	newquery.db = db
	newquery.entities = entities

	// Get the underlying type of the elements of the 'entities' slice
	// cast to interface and save on the query.
	// This will be used for unmarhsalling
	newquery.entity = reflect.New(value.Type().Elem()).Interface()

	return &newquery
}

func (q *query) From(t time.Time) *query {
	q.from = gouuidv6.NewFromTime(t)
	return q
}

func (q *query) To(t time.Time) *query {
	q.to = gouuidv6.NewFromTime(t)
	return q
}

func (q *query) Limit(n int) *query {
	q.limit = n
	return q
}

func (q *query) Reverse() *query {
	q.reverse = true
	return q
}

func (q query) rangePrefixes() (from, to []byte) {
	if !q.from.IsNil() {
		from = makePrefix(q.keyRoot, q.from.Bytes())
	} else {
		from = makePrefix(q.keyRoot, []byte{})
	}

	to = makePrefix(q.keyRoot, []byte{})

	return
}

func (q *query) Run() (int, error) {
	options := badger.DefaultIteratorOptions
	options.Reverse = q.reverse

	return q.execute(options, true)
}

func (q *query) Count() (int, error) {
	options := badger.DefaultIteratorOptions
	options.Reverse = q.reverse
	options.PrefetchValues = false

	return q.execute(options, false)
}

func (q *query) execute(options badger.IteratorOptions, getValues bool) (int, error) {
	// Work out what prefix to iterate over
	from, to := q.rangePrefixes()

	compareKey := compareKey(q.to, to)

	// Cast the 'entity' we have stored on the query to a 'Tormentable'
	// so that it can be unmarshalled
	entity := q.entity.(Tormentable)

	// Set up a slice to accumulate the results
	results := []Tormentable{}
	counter := 0

	err := q.db.KV.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(options)
		defer it.Close()

		for it.Seek(from); it.ValidForPrefix(to); it.Next() {
			if q.limit > 0 && counter >= q.limit {
				break
			}

			key := it.Item().Key()
			if !q.to.IsNil() && !compare(compareKey, key, q.reverse) {
				break
			}

			// Only unmarshall and append to results
			// If we are running a full query, not a count
			if getValues {
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

			// For counts, instead of appending all the results and taking the length
			// we just use a simple counter
			counter++
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	if getValues {
		// Now we have a slice of 'Tormentables'

		// Set up a slice for 'translating' the 'Tormentables' into the target slice type
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

	return counter, nil

}
