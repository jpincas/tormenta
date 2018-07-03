package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta"
)

const noOrders = 1000000

func Benchmark_Count1M(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var ordersToSave []tormenta.Tormentable

	for i := 0; i < noOrders; i++ {
		ordersToSave = append(ordersToSave, &tormenta.Order{
			Customer: i,
		})
	}

	db.Save(ordersToSave...)

	var orders []tormenta.Order

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.Query(&orders).Count()
	}
}

func Benchmark_QueryRun1M(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var ordersToSave []tormenta.Tormentable

	for i := 0; i < noOrders; i++ {
		ordersToSave = append(ordersToSave, &tormenta.Order{
			Customer: i,
		})
	}

	db.Save(ordersToSave...)

	var orders []tormenta.Order

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.Query(&orders).Run()
	}
}
