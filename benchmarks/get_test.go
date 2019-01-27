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
