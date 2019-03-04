package tormenta_test

import (
	"encoding/json"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/tormenta"
	jsoniter "github.com/json-iterator/go"
	"github.com/pquerna/ffjson/ffjson"
)

// All tests use testDBOptions to open the DB
// Just commment out, leaving the set of options you want to use to run the tests

var testDBOptions tormenta.Options = testOptionsStdLib

// var testDBOptions tormenta.Options = testOptionsFFJson
// var testDBOptions tormenta.Options = testOptionsJSONIterFastest
// var testDBOptions tormenta.Options = testOptionsJSONIterDefault
// var testDBOptions tormenta.Options = testOptionsJSONIterCompatible

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
