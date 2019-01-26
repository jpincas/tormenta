package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta"
)

const noOrders = 1000

func Benchmark_Count1M(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var ordersToSave []tormenta.Record

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
		db.query(&orders).Count()
	}
}

func Benchmark_queryRun1M(b *testing.B) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var ordersToSave []tormenta.Record

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
		db.query(&orders).Run()
	}
}
