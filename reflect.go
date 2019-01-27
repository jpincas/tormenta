package tormenta

import "reflect"

// The idea here is to keep all the reflect code in one place,
// which might help to spot potential optimisations / refactors

func newRecord(sliceTarget interface{}) Record {
	_, value := entityTypeAndValue(sliceTarget)
	typ := value.Type().Elem()
	return reflect.New(typ).Interface().(Record)
}

func newResultsArray(sliceTarget interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(sliceTarget))
}

func recordValue(record Record) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(record))
}
