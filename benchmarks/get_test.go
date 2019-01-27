package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Benchmark_QueryGet(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	toSave := stdRecord()
	db.Save(toSave)
	id := toSave.GetID()

	// Reuse the same results
	result := testtypes.FullStruct{}

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.Get(&result, id)
	}
}

func Benchmark_QueryGetIDs(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	toSave1 := stdRecord()
	toSave2 := stdRecord()
	toSave3 := stdRecord()
	toSave4 := stdRecord()
	toSave5 := stdRecord()

	db.Save(toSave1, toSave2, toSave2, toSave4, toSave5)

	// Reuse the same results
	result := []testtypes.FullStruct{}

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.GetIDs(&result, toSave1.GetID(), toSave2.GetID(), toSave3.GetID(), toSave4.GetID(), toSave5.GetID())
	}
}
