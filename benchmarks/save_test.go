package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta"
	jsoniter "github.com/json-iterator/go"
)

func Benchmark_Save_JSONIter_Fastest(b *testing.B) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", tormenta.Options{
		Serialiser:  tormenta.SerialiserJSONIter,
		JsonIterAPI: jsoniter.ConfigFastest,
	})
	defer db.Close()

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.Save(stdRecord())
	}
}

func Benchmark_Save_FFJson(b *testing.B) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", tormenta.Options{
		Serialiser: tormenta.SerialiserJSONff,
	})
	defer db.Close()

	var toSave []tormenta.Record

	for i := 0; i < nRecords; i++ {
		toSave = append(toSave, stdRecord())
	}

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.Save(stdRecord())
	}
}

func Benchmark_Save_StdLib(b *testing.B) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", tormenta.Options{
		Serialiser: tormenta.SerialiserJSONStdLib,
	})
	defer db.Close()

	var toSave []tormenta.Record

	for i := 0; i < nRecords; i++ {
		toSave = append(toSave, stdRecord())
	}

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.Save(stdRecord())
	}
}

func Benchmark_SaveIndividually(b *testing.B) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	var toSave []tormenta.Record

	for i := 0; i < nRecords; i++ {
		toSave = append(toSave, stdRecord())
	}

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.SaveIndividually(toSave...)
	}
}
