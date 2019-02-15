package tormenta_test

import (
	"fmt"
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

// Test range queries across different types
func Test_OrderBy(t *testing.T) {
	var fullStructs []tormenta.Record

	for i := 0; i < 10; i++ {
		// Notice that the intField and StringField increment in oposite ways,
		// such that sorting by either field will produce inverse results.
		// Also - we only go up to 9 so as to avoid alphabetical sorting
		// issues with numbers prefixed by 0
		fullStructs = append(fullStructs, &testtypes.FullStruct{
			IntField:    10 - i,
			StringField: fmt.Sprintf("int-%v", i),
		})
	}

	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()
	db.Save(fullStructs...)

	// First try ordering by intField

	intFieldResults := []testtypes.FullStruct{}
	n, err := db.Find(&intFieldResults).OrderBy("intfield").Run()

	if err != nil {
		t.Errorf("Testing ORDER BY intfield, got error %s", err)
	}

	if n != len(fullStructs) {
		t.Errorf("Testing ORDER BY intfield, n (%v) does not equal actual number of records saved (%v)", n, len(fullStructs))
	}

	if n != len(intFieldResults) {
		t.Errorf("Testing ORDER BY intfield, n (%v) does not equal actual number of results (%v)", n, len(intFieldResults))
	}

	if intFieldResults[0].IntField != 1 {
		t.Errorf("Testing ORDER BY intfield, first member should be 1 but is %v", intFieldResults[0].IntField)
	}

	if intFieldResults[len(intFieldResults)-1].IntField != 10 {
		t.Errorf("Testing ORDER BY intfield, last member should be 10 but is %v", intFieldResults[len(intFieldResults)-1].IntField)
	}

	// First try ordering by stringField

	stringFieldResults := []testtypes.FullStruct{}
	n, err = db.Find(&stringFieldResults).OrderBy("stringfield").Run()

	if err != nil {
		t.Errorf("Testing ORDER BY stringfield, got error %s", err)
	}

	if n != len(fullStructs) {
		t.Errorf("Testing ORDER BY stringfield, n (%v) does not equal actual number of records saved (%v)", n, len(fullStructs))
	}

	if n != len(stringFieldResults) {
		t.Errorf("Testing ORDER BY stringfield, n (%v) does not equal actual number of results (%v)", n, len(stringFieldResults))
	}

	if stringFieldResults[0].StringField != "int-0" {
		t.Errorf("Testing ORDER BY stringfield, first member should be int-0 but is %s", stringFieldResults[0].StringField)
	}

	if stringFieldResults[len(stringFieldResults)-1].StringField != "int-9" {
		t.Errorf("Testing ORDER BY stringfield, last member should be int-9 but is %s", stringFieldResults[len(intFieldResults)-1].StringField)
	}

	// Now compare first members and make sure they are different
	if intFieldResults[0].ID == stringFieldResults[0].ID {
		t.Errorf("Testing ORDER BY. ID's of first member of both results arrays are the same")
	}

	// Now compare last members and make sure they are different
	if intFieldResults[len(intFieldResults)-1].ID == stringFieldResults[len(stringFieldResults)-1].ID {
		t.Errorf("Testing ORDER BY. ID's of first member of both results arrays are the same")
	}

	//Now compare first and last members and make sure they are the same
	if intFieldResults[0].ID != stringFieldResults[len(stringFieldResults)-1].ID {
		t.Errorf("Testing ORDER BY.  First member of array A should be the same as last member of Array B but got %v vs %v", intFieldResults[0].IntField, stringFieldResults[len(stringFieldResults)-1].IntField)
	}
}
