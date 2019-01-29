package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func TestRandomise(t *testing.T) {
	// Make a list of 100 fullStructs
	var fullStructs []tormenta.Record
	for i := 0; i <= 100; i++ {
		fullStructs = append(fullStructs, &testtypes.FullStruct{IntField: i})
	}

	// Make a copy of the list before randomising, then randomise
	ttsBeforeRand := make([]tormenta.Record, len(fullStructs))
	copy(ttsBeforeRand, fullStructs)
	tormenta.RandomiseRecords(fullStructs)

	// Go through element by element, compare, and set a flag to true if a difference was found
	foundDiff := false
	for i := range fullStructs {
		if fullStructs[i].(*testtypes.FullStruct).IntField != ttsBeforeRand[i].(*testtypes.FullStruct).IntField {
			foundDiff = true
		}
	}

	// If no differences were found, then fail
	if !foundDiff {
		t.Error("Testing randomise slice. Could not find any differences after randomisation")
	}

}
