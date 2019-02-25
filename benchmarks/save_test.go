package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta"
)

func Benchmark_Save_JSONIter_Fastest(b *testing.B) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testOptionsJSONIterFastest)
	defer db.Close()

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.Save(stdRecord())
	}
}

func Benchmark_Save_FFJson(b *testing.B) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testOptionsFFJson)
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
	db, _ := tormenta.OpenTestWithOptions("data/tests", testOptionsStdLib)
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
