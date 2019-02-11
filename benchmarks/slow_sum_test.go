package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Benchmark_SlowSum_Test(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	var toSave []tormenta.Record

	n := 1000
	for i := 0; i < n; i++ {
		toSave = append(toSave, stdRecord())
	}

	db.SaveIndividually(toSave...)
	var results []testtypes.FullStruct
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		db.Find(&results).Sum([]string{"IntField"})
	}
}
