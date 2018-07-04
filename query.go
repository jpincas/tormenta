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
	holder   interface{}
	target   interface{}
	from, to gouuidv6.UUID
	limit    int
	reverse  bool
	first    bool
}

func (db DB) newQuery(target interface{}, first bool) *query {
	// Get the key root and cache the value
	keyRoot, value := getKeyRoot(target)

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

	q := &query{
		db:      db,
		keyRoot: keyRoot,
		value:   value,
		holder:  holder,
		target:  target,
	}

	if first {
		q.limit = 1
		q.first = true
	}

	return q
}

// Find kicks off a DB query returning a slice of entities
func (db DB) Find(entities interface{}) *query {
	return db.newQuery(entities, false)
}

// First kicks off a DB query returning the first entity that matches the criteria
func (db DB) First(entity interface{}) *query {
	return db.newQuery(entity, true)
}

// From adds a lower boundary to the date range of the query
func (q *query) From(t time.Time) *query {
	q.from = gouuidv6.NewFromTime(t)
	return q
}

// To adds an upper bound to the date range of the query
func (q *query) To(t time.Time) *query {
	q.to = gouuidv6.NewFromTime(t)
	return q
}

// Limit limits the number of results a query will return to n
func (q *query) Limit(n int) *query {
	q.limit = n
	return q
}

// Reverse reverses the order of date range scanning and returned results (i.e. scans from 'new' to 'old', instead of the default 'old' to 'new' )
func (q *query) Reverse() *query {
	q.reverse = true
	return q
}

// Run actually executes the query and returns results
func (q *query) Run() (int, error) {
	return q.execute(true)
}

// Count executes the query in fast, count-only mode
func (q *query) Count() (int, error) {
	return q.execute(false)
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

func (q query) getIteratorOptions(getValues bool) badger.IteratorOptions {
	options := badger.DefaultIteratorOptions
	options.Reverse = q.reverse
	options.PrefetchValues = getValues
	return options
}

func (q *query) execute(getValues bool) (int, error) {
	// Get the from and to prefix for matching in the iterator
	from, to := q.rangePrefixes()

	// Get the 'to' key for comparison
	compareKey := compareKey(q.to, to)

	// Cast the 'entity' we have stored on the query to a 'Tormentable'
	// so that it can be unmarshalled
	entity := q.holder.(Tormentable)

	// Set up a slice to accumulate the results
	results := []Tormentable{}

	// Initialise the counter, mainly for the count operation
	counter := 0

	err := q.db.KV.View(func(txn *badger.Txn) error {

		// Initialise the iterator, either with or without value pre-fetch
		it := txn.NewIterator(q.getIteratorOptions(getValues))
		defer it.Close()

		// Start iteration
		for it.Seek(from); it.ValidForPrefix(to); it.Next() {

			// If a limit has been specifed AND
			// the counter has reached that limit,
			// then break out by returning 'nil' from the iteration
			if q.limit > 0 && counter >= q.limit {
				return nil
			}

			// If a 'to' clause has been added to the range
			// then compare the current key of the iterator
			// to the theoretical final key in the range
			// and break out if we've reached it
			key := it.Item().Key()
			if !q.to.IsNil() && !compare(compareKey, key, q.reverse) {
				return nil
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

				// If this is a 'first only' query
				// then we can break out at this point
				// Set the counter to 1 to show we actually found it
				if q.first {
					counter = 1
					return nil
				}

				// Otherwise we append the unmarshalled result
				// and carry on iterating
				results = append(results, entity)
			}

			// For counts, instead of appending all the results and taking the length
			// we just use a simple counter
			counter++
		}

		// End the iteration
		return nil
	})

	// If there was an error from the iterator
	if err != nil {
		return 0, err
	}

	// For a first only query,
	// we just need to set the value of the single result
	// If nothing was found, then the counter will not have been set to 1
	// so it will still be 0
	if q.first {
		reflect.Indirect(reflect.ValueOf(q.target)).Set(reflect.Indirect(reflect.ValueOf(entity)))
		return counter, nil
	}

	if getValues {

		// Now we have a slice of 'Tormentables'

		// Set up a slice for 'translating' the 'Tormentables' into the target slice type
		rt := reflect.Indirect(reflect.ValueOf(q.target))

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
		reflect.Indirect(reflect.ValueOf(q.target)).Set(rt)

		return len(results), nil
	}

	return counter, nil

}
