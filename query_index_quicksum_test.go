package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

// Helper for making groups of depatments
func getDept(i int) int {
	if i <= 10 {
		return 1
	} else if i <= 20 {
		return 2
	} else {
		return 3
	}
}

// Test aggregation on an index
func Test_QuickSum(t *testing.T) {
	var fullStructs []tormenta.Record

	// Accumulators
	var accInt int
	var accInt16 int16
	var accInt32 int32
	var accInt64 int64

	var accUint uint
	var accUint16 uint16
	var accUint32 uint32
	var accUint64 uint64

	var accFloat32 float32
	var accFloat64 float64

	// Range - assymetric neg/pos so total doesn't balance out
	for i := -30; i <= 100; i++ {

		fullStruct := &testtypes.FullStruct{
			// Signed Ints
			IntField:   i,
			Int16Field: int16(i),
			Int32Field: int32(i),
			Int64Field: int64(i),

			// Unsigned Ints
			UintField:   uint(i * i),
			Uint16Field: uint16(i * i),
			Uint32Field: uint32(i * i),
			Uint64Field: uint64(i * i),

			// Floats
			FloatField:   float64(i),
			Float32Field: float32(i),
		}

		accInt += i
		accInt16 += int16(i)
		accInt32 += int32(i)
		accInt64 += int64(i)

		accUint += uint(i * i)
		accUint16 += uint16(i * i)
		accUint32 += uint32(i * i)
		accUint64 += uint64(i * i)

		accFloat64 += float64(i)
		accFloat32 += float32(i)

		fullStructs = append(fullStructs, fullStruct)
	}

	// Randomise and save
	tormenta.RandomiseRecords(fullStructs)
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()
	db.Save(fullStructs...)

	var resultInt int
	var resultInt16 int16
	var resultInt32 int32
	var resultInt64 int64

	var resultUint uint
	var resultUint16 uint16
	var resultUint32 uint32
	var resultUint64 uint64

	var resultFloat32 float32
	var resultFloat64 float64

	// Test cases
	testCases := []struct {
		name      string
		fieldName string
		sumResult interface{}
		acc       interface{}
		// Specify how to convert back the results pointer into a comparable value
		convertBack func(interface{}) interface{}
	}{
		// Ints
		{"int", "intfield", &resultInt, accInt, func(n interface{}) interface{} { return *n.(*int) }},
		{"int16", "int16field", &resultInt16, accInt16, func(n interface{}) interface{} { return *n.(*int16) }},
		{"int32", "int32field", &resultInt32, accInt32, func(n interface{}) interface{} { return *n.(*int32) }},
		{"int64", "int64field", &resultInt64, accInt64, func(n interface{}) interface{} { return *n.(*int64) }},

		// Uints
		{"uint", "uintfield", &resultUint, accUint, func(n interface{}) interface{} { return *n.(*uint) }},
		{"uint16", "uint16field", &resultUint16, accUint16, func(n interface{}) interface{} { return *n.(*uint16) }},
		{"uint32", "uint32field", &resultUint32, accUint32, func(n interface{}) interface{} { return *n.(*uint32) }},
		{"uint64", "uint64field", &resultUint64, accUint64, func(n interface{}) interface{} { return *n.(*uint64) }},

		// Floats
		{"float32", "float32field", &resultFloat32, accFloat32, func(n interface{}) interface{} { return *n.(*float32) }},
		{"float64", "floatfield", &resultFloat64, accFloat64, func(n interface{}) interface{} { return *n.(*float64) }},
	}

	for _, test := range testCases {
		results := []testtypes.FullStruct{}

		// We use OrderBy in order to specify an index without any range -
		// see notes on OrderBy
		_, err := db.Find(&results).OrderBy(test.fieldName).QuickSum(test.sumResult)

		// Check for error
		if err != nil {
			t.Errorf("Testing %s quicksum.  Got error: %s", test.name, err)
		}

		// Compare result to accumulator
		result := test.convertBack(test.sumResult)
		if result != test.acc {
			t.Errorf("Testing %s quicksum. Expected %v, got %v", test.name, test.acc, result)
		}
	}
}
