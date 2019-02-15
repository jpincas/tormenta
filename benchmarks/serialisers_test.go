package benchmarks

import (
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
	jsoniter "github.com/json-iterator/go"
)

func prepareDB(dbOptions tormenta.Options) {
	noEntities := 1000

	// Open the DB
	db, err := tormenta.OpenTest("data/tests", dbOptions)
	if err != nil {
		panic("failed to open db")
	}
	defer db.Close()

	entity := testtypes.FullStruct{
		// Basic Types
		IntField:    1,
		StringField: "test",
		FloatField:  0.99,
		BoolField:   true,
		// Note: time.Now() includes a monotonic clock component, which is stripped
		// for marshalling, which destroys equality between saved and retrieved, even
		// though they are essentially the 'same'. See: https://golang.org/pkg/time/
		// We therefore use a fixed time
		DateField: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),

		// Slice Types
		IDSliceField: []gouuidv6.UUID{
			gouuidv6.New(),
			gouuidv6.New(),
			gouuidv6.New(),
		},
		IntSliceField:    []int{1, 2},
		StringSliceField: []string{"test1", "test2"},
		FloatSliceField:  []float64{0.99, 1.99},
		BoolSliceField:   []bool{true, false},
		DateSliceField: []time.Time{
			time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC),
			time.Date(2011, time.November, 10, 23, 0, 0, 0, time.UTC),
		},

		// Basic Defined Fields
		DefinedIntField:    testtypes.DefinedInt(1),
		DefinedStringField: testtypes.DefinedString("test"),
		DefinedFloatField:  testtypes.DefinedFloat(0.99),
		DefinedBoolField:   testtypes.DefinedBool(true),
		DefinedIDField:     testtypes.DefinedID(gouuidv6.New()),
		DefinedDateField:   testtypes.DefinedDate(time.Date(2011, time.November, 10, 23, 0, 0, 0, time.UTC)),

		// Defined Slice Type Fields
		DefinedIDSliceField: []testtypes.DefinedID{
			testtypes.DefinedID(gouuidv6.New()),
			testtypes.DefinedID(gouuidv6.New()),
			testtypes.DefinedID(gouuidv6.New()),
		},
		DefinedIntSliceField:    []testtypes.DefinedInt{1, 2},
		DefinedStringSliceField: []testtypes.DefinedString{"test1", "test2"},
		DefinedFloatSliceField:  []testtypes.DefinedFloat{0.99, 1.99},
		DefinedBoolSliceField:   []testtypes.DefinedBool{true, false},
		DefinedDateSliceField: []testtypes.DefinedDate{
			testtypes.DefinedDate(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)),
			testtypes.DefinedDate(time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)),
			testtypes.DefinedDate(time.Date(2011, time.November, 10, 23, 0, 0, 0, time.UTC)),
		},

		// Embedded Struct
		MyStruct: testtypes.MyStruct{
			StructIntField:    1,
			StructStringField: "test",
			StructFloatField:  0.99,
			StructBoolField:   true,
			StructDateField:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		},

		// Named struct field
		StructField: testtypes.MyStruct{
			StructIntField:    1,
			StructStringField: "test",
			StructFloatField:  0.99,
			StructBoolField:   true,
			StructDateField:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		},
	}

	for i := 0; i < noEntities; i++ {
		// Save
		entity.SetID(gouuidv6.New())
		db.Save(&entity)
	}

	for i := 0; i < noEntities; i++ {
		// Get
		// var result testtypes.FullStruct
		// db.Get(&result, entity.ID)
	}

}

func Benchmark_Serialisers_JsonIter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		prepareDB(tormenta.Options{
			JsonIterAPI: jsoniter.ConfigFastest,
			Serialiser:  tormenta.SerialiserJSONIter,
		})
	}
}

func Benchmark_Serialisers_StdLib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		prepareDB(tormenta.Options{
			Serialiser: tormenta.SerialiserJSONStdLib,
		})
	}
}
