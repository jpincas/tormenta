package tormenta_test

import (
	"testing"
	"time"

	"github.com/jpincas/tormenta"
)

func Test_BasicQuery(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// 1 tt
	tt1 := TestType{}
	db.Save(&tt1)

	var tts []TestType
	n, _, err := db.Find(&tts).Run()

	if err != nil {
		t.Error("Testing basic querying - got error")
	}

	if len(tts) != 1 || n != 1 {
		t.Errorf("Testing querying with 1 entity saved. Expecting 1 entity - got %v/%v", len(tts), n)
	}

	tts = []TestType{}
	c, _, err := db.Find(&tts).Count()
	if c != 1 {
		t.Errorf("Testing count 1 entity saved. Expecting 1 - got %v", c)
	}

	// 2 tts
	tt2 := TestType{}
	db.Save(&tt2)

	tts = []TestType{}
	if n, _, _ := db.Find(&tts).Run(); n != 2 {
		t.Errorf("Testing querying with 2 entity saved. Expecting 2 entities - got %v", n)
	}

	if c, _, _ := db.Find(&tts).Count(); c != 2 {
		t.Errorf("Testing count 2 entities saved. Expecting 2 - got %v", c)
	}
	if tt1.ID == tt2.ID {
		t.Errorf("Testing querying with 2 entities saved. 2 entities saved both have same ID")
	}
	if tts[0].ID == tts[1].ID {
		t.Errorf("Testing querying with 2 entities saved. 2 results returned. Both have same ID")
	}

	// Limit
	tts = []TestType{}
	if n, _, _ := db.Find(&tts).Limit(1).Run(); n != 1 {
		t.Errorf("Testing querying with 2 entities saved + limit. Wrong number of results received")
	}

	// Reverse - simple, only tests number received
	tts = []TestType{}
	if n, _, _ := db.Find(&tts).Reverse().Run(); n != 2 {
		t.Errorf("Testing querying with 2 entities saved + reverse. Expected %v, got %v", 2, n)
	}

	// Reverse + Limit - simple, only tests number received
	tts = []TestType{}
	if n, _, _ := db.Find(&tts).Reverse().Limit(1).Run(); n != 1 {
		t.Errorf("Testing querying with 2 entities saved + reverse + limit. Expected %v, got %v", 1, n)
	}

	// Reverse + Count
	tts = []TestType{}
	if c, _, _ := db.Find(&tts).Reverse().Count(); c != 2 {
		t.Errorf("Testing count with 2 entities saved + reverse. Expected %v, got %v", 2, c)
	}

}

func Test_BasicQuery_First(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	tt1 := TestType{}
	tt2 := TestType{}
	db.Save(&tt1, &tt2)

	var tt TestType
	n, _, err := db.First(&tt).Run()

	if err != nil {
		t.Error("Testing first - got error")
	}

	if n != 1 {
		t.Errorf("Testing first. Expecting 1 entity - got %v", n)
	}

	if tt.ID.IsNil() {
		t.Errorf("Testing first. Nil ID retrieved")
	}

	if tt.ID != tt1.ID {
		t.Errorf("Testing first. Order IDs are not equal - wrong tt retrieved")
	}

	// Test nothing found (impossible range)
	n, _, _ = db.First(&tt).From(time.Now()).To(time.Now()).Run()
	if n != 0 {
		t.Errorf("Testing first when nothing should be found.  Got n = %v", n)
	}
}
