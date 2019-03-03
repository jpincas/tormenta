package benchmarks

import (
	"encoding/json"

	"github.com/dgraph-io/badger"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
	jsoniter "github.com/json-iterator/go"
	"github.com/pquerna/ffjson/ffjson"
)

const nRecords = 1000

func stdRecord() *testtypes.FullStruct {
	return &testtypes.FullStruct{
		IntField:          1,
		StringField:       "test",
		MultipleWordField: "multiple word field",
		FloatField:        9.99,
		BoolField:         true,
		IntSliceField:     []int{1, 2, 3, 4, 5},
		StringSliceField:  []string{"string", "slice", "field"},
		FloatSliceField:   []float64{0.1, 0.2, 0.3, 0.4, 0.5},
		BoolSliceField:    []bool{true, false, true, false},
		MyStruct: testtypes.MyStruct{
			StructIntField:    100,
			StructFloatField:  999.999,
			StructBoolField:   false,
			StructStringField: "embedded string field",
		},
	}
}

var testDBOptions = testOptionsStdLib

// var testDBOptions = testOptionsFFJson
// var testDBOptions = testOptionsJSONIterFastest
// var testDBOptions = testOptionsJSONIterDefault
// var testDBOptions = testOptionsJSONIterCompatible

var testOptionsStdLib = tormenta.Options{
	SerialiseFunc:   json.Marshal,
	UnserialiseFunc: json.Unmarshal,
	BadgerOptions:   badger.DefaultOptions,
}

var testOptionsFFJson = tormenta.Options{
	SerialiseFunc:   ffjson.Marshal,
	UnserialiseFunc: ffjson.Unmarshal,
	BadgerOptions:   badger.DefaultOptions,
}

var testOptionsJSONIterFastest = tormenta.Options{
	// Main difference is precision of floats - see https://godoc.org/github.com/json-iterator/go
	SerialiseFunc:   jsoniter.ConfigFastest.Marshal,
	UnserialiseFunc: jsoniter.ConfigFastest.Unmarshal,
	BadgerOptions:   badger.DefaultOptions,
}

var testOptionsJSONIterDefault = tormenta.Options{
	// Main difference is precision of floats - see https://godoc.org/github.com/json-iterator/go
	SerialiseFunc:   jsoniter.ConfigDefault.Marshal,
	UnserialiseFunc: jsoniter.ConfigDefault.Unmarshal,
	BadgerOptions:   badger.DefaultOptions,
}

var testOptionsJSONIterCompatible = tormenta.Options{
	// Main difference is precision of floats - see https://godoc.org/github.com/json-iterator/go
	SerialiseFunc:   jsoniter.ConfigCompatibleWithStandardLibrary.Marshal,
	UnserialiseFunc: jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal,
	BadgerOptions:   badger.DefaultOptions,
}
