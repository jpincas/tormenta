package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
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
func Test_Aggregation(t *testing.T) {
	var tts []tormenta.Record

	for i := 1; i <= 30; i++ {
		tt := &TestType{
			FloatField: float64(i),
			IntField:   i,
		}

		tts = append(tts, tt)
	}

	tormenta.RandomiseRecords(tts)

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(tts...)

	results := []TestType{}
	var intSum int32
	var floatSum float64
	expected := 465

	// Int32

	_, _, err := db.Find(&results).Range("intfield", 1, 30).Sum(&intSum)
	if err != nil {
		t.Error("Testing int32 agreggation.  Got error")
	}

	expectedIntSum := int32(expected)
	if intSum != expectedIntSum {
		t.Errorf("Testing int32 agreggation. Expteced %v, got %v", expectedIntSum, intSum)
	}

	// Float64

	_, _, err = db.Find(&results).Range("floatfield", 1.00, 30.00).Sum(&floatSum)
	if err != nil {
		t.Error("Testing float64 agreggation.  Got error")
	}

	expectedFloatSum := float64(expected)
	if floatSum != expectedFloatSum {
		t.Errorf("Testing float64 agreggation. Expteced %v, got %v", expectedFloatSum, floatSum)
	}
}
