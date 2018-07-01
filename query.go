package tormenta

import "reflect"

type Query struct {
	keyRoot  []byte
	value    reflect.Value
	entities interface{}
}

func (db DB) Query(entities interface{}) Query {
	newQuery := Query{}
	keyRoot, value := getKeyRoot(entities)
	newQuery.keyRoot = keyRoot
	newQuery.value = value
	newQuery.entities = entities

	return newQuery
}

func (q Query) Run() error {
	orders := []Order{
		{Model: newModel()},
		{Model: newModel()},
		{Model: newModel()},
	}

	reflect.
		Indirect(q.value).
		Set(reflect.ValueOf(orders))

	return nil
}
