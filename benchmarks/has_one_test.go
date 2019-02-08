package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta/testtypes"

	"github.com/jpincas/tormenta"
)

func Benchmark_Relations_HasOne(b *testing.B) {
	noEntities := 1000
	noRelations := 10

	// Open the DB
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	// Created some nested structs and save
	// nestedStruct1 := testtypes.NestedRelatedStruct{}
	// nestedStruct2 := testtypes.NestedRelatedStruct{}
	// db.Save(&nestedStruct1, &nestedStruct2)

	// Create some related structs which nest the above and save
	var relatedStructs []tormenta.Record
	for i := 0; i < noRelations; i++ {
		relatedStruct := testtypes.RelatedStruct{}
		relatedStructs = append(relatedStructs, &relatedStruct)
	}
	db.Save(relatedStructs...)

	// Create some full structs including these relations
	// To make things a little more realistic, we will rotate relations,
	// repeated N relations using %
	var fullStructs []tormenta.Record
	for i := 0; i < noEntities; i++ {
		fullStruct := testtypes.FullStruct{
			HasOneID:        relatedStructs[i%noRelations].GetID(),
			HasAnotherOneID: relatedStructs[i%noRelations].GetID(),
		}

		fullStructs = append(fullStructs, &fullStruct)
	}
	db.Save(fullStructs...)

	// Reset the timer
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tormenta.HasOne(db, []string{"HasOne", "HasAnotherOne"}, fullStructs...)
	}

}
