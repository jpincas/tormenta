package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
)

func TestRandomise(t *testing.T) {
	// Make a list of 100 tts
	var tts []tormenta.Record
	for i := 0; i <= 100; i++ {
		tts = append(tts, &TestType{IntField: i})
	}

	// Make a copy of the list before randomising, then randomise
	ttsBeforeRand := make([]tormenta.Record, len(tts))
	copy(ttsBeforeRand, tts)
	tormenta.RandomiseRecords(tts)

	// Go through element by element, compare, and set a flag to true if a difference was found
	foundDiff := false
	for i := range tts {
		if tts[i].(*TestType).IntField != ttsBeforeRand[i].(*TestType).IntField {
			foundDiff = true
		}
	}

	// If no differences were found, then fail
	if !foundDiff {
		t.Error("Testing randomise slice. Could not find any differences after randomisation")
	}

}
