package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Benchmark_QueryCount(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	var toSave []tormenta.Record

	for i := 0; i < nRecords; i++ {
		toSave = append(toSave, stdRecord())
	}

	db.Save(toSave...)

	var fullStructs []testtypes.FullStruct

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.Find(&fullStructs).Count()
	}
}

func Benchmark_QueryRun(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	var toSave []tormenta.Record

	for i := 0; i < nRecords; i++ {
		toSave = append(toSave, stdRecord())
	}

	db.Save(toSave...)

	var fullStructs []testtypes.FullStruct

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.Find(&fullStructs).Run()
	}
}
