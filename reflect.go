package tormenta

import "reflect"

// The idea here is to keep all the reflect code in one place,
// which might help to spot potential optimisations / refactors

func newRecordFromSlice(target interface{}) Record {
	_, value := entityTypeAndValue(target)
	typ := value.Type().Elem()
	return reflect.New(typ).Interface().(Record)
}

func newRecord(target interface{}) Record {
	_, value := entityTypeAndValue(target)
	typ := value.Type()
	return reflect.New(typ).Interface().(Record)
}

func newResultsArray(sliceTarget interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(sliceTarget))
}

func recordValue(record Record) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(record))
}

func setResultsArrayOntoTarget(sliceTarget interface{}, records reflect.Value) {
	reflect.Indirect(reflect.ValueOf(sliceTarget)).Set(records)
}

func setSingleResultOntoTarget(target interface{}, record Record) {
	reflect.Indirect(reflect.ValueOf(target)).Set(reflect.Indirect(reflect.ValueOf(record)))
}
