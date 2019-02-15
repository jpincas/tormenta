package tormenta_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Test_Save_Get(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	// Test Entity

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
		DefinedStructField: testtypes.DefinedStruct(
			testtypes.MyStruct{
				StructIntField:    1,
				StructStringField: "test",
				StructFloatField:  0.99,
				StructBoolField:   true,
				StructDateField:   time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			},
		),

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

	// Save
	db.Save(&entity)

	// Get
	var result testtypes.FullStruct
	db.Get(&result, entity.ID)

	testCases := []struct {
		name             string
		saved, retrieved interface{}
		deep             bool
	}{
		// Basic Types
		{"ID", entity.ID, result.ID, false},
		{"IntField", entity.IntField, result.IntField, false},
		{"StringField", entity.StringField, result.StringField, false},
		{"FloatField", entity.FloatField, result.FloatField, false},
		{"BoolField", entity.BoolField, result.BoolField, false},
		{"DateField", entity.DateField, result.DateField, false},

		// Slice Types
		{"IDSliceField", entity.IDSliceField, result.IDSliceField, true},
		{"IntSliceField", entity.IntSliceField, result.IntSliceField, true},
		{"StringSliceField", entity.StringSliceField, result.StringSliceField, true},
		{"FloatSliceField", entity.FloatSliceField, result.FloatSliceField, true},
		{"BoolSliceField", entity.BoolSliceField, result.BoolSliceField, true},
		{"DateSliceField", entity.DateSliceField, result.DateSliceField, true},

		// Basic Defined Types
		{"DefinedID", entity.DefinedIDField, result.DefinedIDField, false},
		{"DefinedIntField", entity.DefinedIntField, result.DefinedIntField, false},
		{"DefinedStringField", entity.DefinedStringField, result.DefinedStringField, false},
		{"DefinedFloatField", entity.DefinedFloatField, result.DefinedFloatField, false},
		{"DefinedBoolField", entity.DefinedBoolField, result.DefinedBoolField, false},
		{"DefinedStructField", entity.DefinedStructField, result.DefinedStructField, false},

		// Defined Slice Types
		{"DefinedIDSliceField", entity.DefinedIDSliceField, result.DefinedIDSliceField, true},
		{"DefinedIntSliceField", entity.DefinedIntSliceField, result.DefinedIntSliceField, true},
		{"DefinedStringSliceField", entity.DefinedStringSliceField, result.DefinedStringSliceField, true},
		{"DefinedFloatSliceField", entity.DefinedFloatSliceField, result.DefinedFloatSliceField, true},
		{"DefinedBoolSliceField", entity.DefinedBoolSliceField, result.DefinedBoolSliceField, true},

		// Embedded Struct
		{"StructID", entity.ID, result.ID, false},
		{"StructIntField", entity.StructIntField, result.StructIntField, false},
		{"StructStringField", entity.StructStringField, result.StructStringField, false},
		{"StructFloatField", entity.StructFloatField, result.StructFloatField, false},
		{"StructBoolField", entity.StructBoolField, result.StructBoolField, false},
		{"StructDateField", entity.StructDateField, result.StructDateField, false},

		//Named Struct
		{"EmbeddedStructField", entity.StructField, result.StructField, true},

		// Defined time.Time fields don't serialise - see README
		// These are just here to remind us not to add them again and wonder why they don't work
		// {"DefinedDateField", entity.DefinedDateField, result.DefinedDateField, false},
		// {"DefinedDateSliceField", entity.DefinedDateSliceField, result.DefinedDateSliceField, true},

	}

	for _, test := range testCases[0:] {
		if !test.deep {
			if test.retrieved != test.saved {
				t.Errorf("Testing %s. Equality test failed. Saved = %v; Retrieved = %v", test.name, test.saved, test.retrieved)
			}
		} else {
			if !reflect.DeepEqual(test.retrieved, test.saved) {
				t.Errorf("Testing %s. Deep equality test failed. Saved = %v; Retrieved = %v", test.name, test.saved, test.retrieved)
			}
		}
	}

}
