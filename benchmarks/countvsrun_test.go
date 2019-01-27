package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta"
)

const notts = 1000

func Benchmark_Count1M(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var ttsToSave []tormenta.Record

	for i := 0; i < notts; i++ {
		ttsToSave = append(ttsToSave, &tormenta.Order{
			Customer: i,
		})
	}

	db.Save(ttsToSave...)

	var tts []tormenta.Order

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.query(&tts).Count()
	}
}

func Benchmark_queryRun1M(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var ttsToSave []tormenta.Record

	for i := 0; i < notts; i++ {
		ttsToSave = append(ttsToSave, &tormenta.Order{
			Customer: i,
		})
	}

	db.Save(ttsToSave...)

	var tts []tormenta.Order

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.query(&tts).Run()
	}
}
