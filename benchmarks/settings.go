package benchmarks

import (
	"github.com/jpincas/tormenta/testtypes"
)

const nRecords = 1000

func stdRecord() *testtypes.FullStruct {
	return &testtypes.FullStruct{
		IntField:          1,
		StringField:       "test",
		MultipleWordField: "multiple word field",
		FloatField:        9.99,
		BoolField:         true,
		IntSliceField:     []int{1, 2, 3, 4, 5},
		StringSliceField:  []string{"string", "slice", "field"},
		FloatSliceField:   []float64{0.1, 0.2, 0.3, 0.4, 0.5},
		BoolSliceField:    []bool{true, false, true, false},
		MyStruct: testtypes.MyStruct{
			StructIntField:    100,
			StructFloatField:  999.999,
			StructBoolField:   false,
			StructStringField: "embedded string field",
		},
	}
}
