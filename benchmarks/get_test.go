package benchmarks

import (
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Benchmark_Get(b *testing.B) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
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

func Benchmark_GetIDs(b *testing.B) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	var toSave []tormenta.Record
	var ids []gouuidv6.UUID

	for i := 0; i <= nRecords; i++ {
		id := gouuidv6.NewFromTime(time.Now())
		record := stdRecord()
		record.SetID(id)
		toSave = append(toSave, record)
		ids = append(ids, id)
	}

	db.Save(toSave...)

	// Reuse the same results
	results := []testtypes.FullStruct{}

	// Reset the timer
	b.ResetTimer()

	// Run the aggregation
	for i := 0; i < b.N; i++ {
		db.GetIDs(&results, ids...)
	}
}

// func Benchmark_GetIDsSerial(b *testing.B) {
// 	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
// 	defer db.Close()

// 	var toSave []tormenta.Record
// 	var ids []gouuidv6.UUID

// 	for i := 0; i <= nRecords; i++ {
// 		id := gouuidv6.NewFromTime(time.Now())
// 		record := stdRecord()
// 		record.SetID(id)
// 		toSave = append(toSave, record)
// 		ids = append(ids, id)
// 	}

// 	db.Save(toSave...)

// 	// Reuse the same results
// 	results := []testtypes.FullStruct{}

// 	// Reset the timer
// 	b.ResetTimer()

// 	// Run the aggregation
// 	for i := 0; i < b.N; i++ {
// 		db.GetIDsSerial(&results, ids...)
// 	}
// }
