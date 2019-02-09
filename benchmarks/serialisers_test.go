package benchmarks

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
	jsoniter "github.com/json-iterator/go"
)

func prepareDB(dbOptions tormenta.Options) (*tormenta.DB, []tormenta.Record) {
	noEntities := 1000
	noRelations := 50

	// Open the DB
	db, err := tormenta.OpenTest("data/tests", dbOptions)
	if err != nil {
		panic("failed to open db")
	}

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

	return db, fullStructs
}

func Benchmark_Serialisers_JsonIter(b *testing.B) {
	db, records := prepareDB(tormenta.Options{
		JsonIterAPI: jsoniter.ConfigFastest,
		Serialiser:  tormenta.SerialiserJSONIter,
	})
	defer db.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tormenta.HasOne(
			db,
			[]string{"HasOne.Nested", "HasAnotherOne.Nested"},
			records...,
		)
	}
}

// CAUSES A REFLECT ERROR - don't know why
// func Benchmark_Serialisers_StdLib(b *testing.B) {
// 	db, records := prepareDB(tormenta.Options{
// 		Serialiser: tormenta.SerialiserJSONStdLib,
// 	})
// 	defer db.Close()

// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		tormenta.HasOne(
// 			db,
// 			[]string{"HasOne.Nested", "HasAnotherOne.Nested"},
// 			records...,
// 		)
// 	}
// }
