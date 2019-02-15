package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta/testtypes"

	"github.com/jpincas/tormenta"
)

func Benchmark_Relations_HasOne(b *testing.B) {
	noEntities := 1000
	noRelations := 50

	// Open the DB
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	// Created some nested structs and save
	nestedStruct1 := testtypes.NestedRelatedStruct{}
	db.Save(&nestedStruct1)

	// Create some related structs which nest the above and save
	var relatedStructs []tormenta.Record
	for i := 0; i < noRelations; i++ {
		relatedStruct := testtypes.RelatedStruct{
			NestedID: nestedStruct1.ID,
		}
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
		tormenta.LoadByID(db, []string{"HasOne.Nested", "HasAnotherOne.Nested"}, fullStructs...)
	}

}
