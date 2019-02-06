package tormenta

import "reflect"

var (
	typeInt    = reflect.TypeOf(0)
	typeUint   = reflect.TypeOf(uint(0))
	typeFloat  = reflect.TypeOf(0.99)
	typeString = reflect.TypeOf("")
	typeBool   = reflect.TypeOf(true)
)

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

// convertUnderlying takes an interface and converts its underlying type
// to the target type.  Obviously the underlying types must be convertible
// E.g. NamedInt -> Int
func convertUnderlying(src interface{}, targetType reflect.Type) interface{} {
	return reflect.ValueOf(src).Convert(targetType).Interface()
}

func intInterfaceToInt32(i interface{}) int32 {
	return int32(convertUnderlying(i, typeInt).(int))
}

func intInterfaceToUInt32(i interface{}) uint32 {
	return uint32(convertUnderlying(i, typeUint).(uint))
}
