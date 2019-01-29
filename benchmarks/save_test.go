package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta"
)

func Benchmark_Save(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	var toSave []tormenta.Record

	for i := 0; i < nRecords; i++ {
		toSave = append(toSave, stdRecord())
	}

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.Save(toSave...)
	}
}

func Benchmark_SaveIndividually(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
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
