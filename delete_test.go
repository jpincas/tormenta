package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta"
)

func Test_Delete(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	tt := TestType{}

	db.Save(&tt)

	// Test the the tt has been saved
	retrievedtt := TestType{}
	ok, _, _ := db.Get(&retrievedtt, tt.ID)
	if !ok || tt.ID != retrievedtt.ID {
		t.Error("Testing delete. Test tt not saved correctly")
	}

	// Delete
	n, err := db.Delete("testtype", tt.ID)

	if err != nil {
		t.Errorf("Testing delete. Got error %v", err)
	}

	if n != 1 {
		t.Errorf("Testing delete. Expected n = 1, got n = %v", n)
	}

	// Attempt to retrieve again
	ok, _, _ = db.Get(&retrievedtt, tt.ID)
	if ok {
		t.Error("Testing delete. Supposedly deleted tt found on 2nd get")
	}
}

func Test_Delete_Multiple(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	tt1 := TestType{}
	tt2 := TestType{}
	tt3 := TestType{}

	db.Save(&tt1, &tt2, &tt3)

	// Delete
	n, err := db.Delete("testtype", tt1.ID, tt2.ID, tt3.ID)

	if err != nil {
		t.Errorf("Testing multiple delete. Got error %v", err)
	}

	if n != 3 {
		t.Errorf("Testing multiple delete. Expected n = %v, got n = %v", 3, n)
	}

	var tts []TestType
	c, _, _ := db.Find(&tts).Count()
	if c > 0 {
		t.Errorf("Testing delete. Should have found any tts, but found %v", c)
	}
}
