package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Test_BasicOrQuery(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	var toSave []tormenta.Record
	for i := 0; i < 10; i++ {
		toSave = append(toSave, &testtypes.FullStruct{
			IntField: i,
		})
	}

	db.Save(toSave...)
	results := []testtypes.FullStruct{}

	n, err := tormenta.
		Or(
			db.Find(&results).Match("intfield", 1),
			db.Find(&results).Match("intfield", 2),
			db.Find(&results).Match("intfield", 3),
		).Run()

	if err != nil {
		t.Error("Testing basic OR - got error")
	}

	if n != len(results) {
		t.Errorf("Testing basic OR - n does not equal length of results. N: %v; Length results: %v", n, len(results))
	}

	if n != 3 {
		t.Errorf("Testing basic OR. Wrong number of results. Expected: %v; got: %v", 3, n)
	}
}
