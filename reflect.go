package tormenta

import (
	"fmt"
	"reflect"
)

var (
	typeInt    = reflect.TypeOf(0)
	typeUint   = reflect.TypeOf(uint(0))
	typeFloat  = reflect.TypeOf(0.99)
	typeString = reflect.TypeOf("")
	typeBool   = reflect.TypeOf(true)
)

// The idea here is to keep all the reflect code in one place,
// which might help to spot potential optimisations / refactors

func indexStringForThisEntity(record Record) string {
	return string(typeToIndexString(reflect.TypeOf(record).String()))
}

func entityTypeAndValue(t interface{}) ([]byte, reflect.Value) {
	e := reflect.Indirect(reflect.ValueOf(t))
	return typeToKeyRoot(e.Type().String()), e
}

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

func fieldValue(entity Record, fieldName string) reflect.Value {
	return recordValue(entity).FieldByName(fieldName)
}

func fieldKind(target interface{}, fieldName string) (reflect.Kind, error) {
	// The target will either be a pointer to slice or struct
	// Start of assuming its a pointer to a struct
	ss := reflect.ValueOf(target).Elem()

	// Check if its a slice, and if so,
	// get the underlying type, create a new struct value pointer,
	// and dereference it
	if ss.Type().Kind() == reflect.Slice {
		ss = reflect.New(ss.Type().Elem()).Elem()
	}

	// At this point, independently of whether the input was a struct or slice,
	// we can get the required field by name and get its kind
	v := ss.FieldByName(fieldName)
	if !v.IsValid() {
		return 0, fmt.Errorf(ErrFieldCouldNotBeFound, fieldName)
	}

	return v.Type().Kind(), nil
}

// newSlice sets up a new target slice for results
// this was arrived at after a lot of experimentation
// so might not be the most efficient way!! TODO
func newSlice(t reflect.Type, l int) interface{} {
	asSlice := reflect.MakeSlice(reflect.SliceOf(t), 0, l)
	new := reflect.New(asSlice.Type())
	new.Elem().Set(asSlice)
	return new.Interface()
}
