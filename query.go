package tormenta

import (
	"reflect"

	"github.com/dgraph-io/badger"
)

type Query struct {
	db      DB
	keyRoot []byte
	value   reflect.Value
	entity  tormentable
}

func (db DB) Query(entity tormentable) Query {
	newQuery := Query{}
	keyRoot, value := getKeyRoot(entity)
	newQuery.keyRoot = keyRoot
	newQuery.value = value
	newQuery.entity = entity
	newQuery.db = db

	return newQuery
}

func (q Query) Run() ([]tormentable, error) {

	results := []tormentable{}
	prefix := makePrefix(q.keyRoot)

	err := q.db.KV.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			val, err := it.Item().Value()
			if err != nil {
				return err
			}

			_, err = q.entity.UnmarshalMsg(val)
			if err != nil {
				return err
			}

			results = append(results, q.entity)
		}

		return nil
	})

	if err != nil {
		return []tormentable{}, err
	}

	return results, nil
}
