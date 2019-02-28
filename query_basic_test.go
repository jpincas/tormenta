package tormenta_test

import (
	"testing"
	"time"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

func Test_BasicQuery(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	// 1 fullStruct
	tt1 := testtypes.FullStruct{}
	db.Save(&tt1)

	var fullStructs []testtypes.FullStruct
	n, err := db.Find(&fullStructs).Run()

	if err != nil {
		t.Error("Testing basic querying - got error")
	}

	if len(fullStructs) != 1 || n != 1 {
		t.Errorf("Testing querying with 1 entity saved. Expecting 1 entity - got %v/%v", len(fullStructs), n)
	}

	fullStructs = []testtypes.FullStruct{}
	c, err := db.Find(&fullStructs).Count()
	if c != 1 {
		t.Errorf("Testing count 1 entity saved. Expecting 1 - got %v", c)
	}

	// 2 fullStructs
	tt2 := testtypes.FullStruct{}
	db.Save(&tt2)
	if tt1.ID == tt2.ID {
		t.Errorf("Testing querying with 2 entities saved. 2 entities saved both have same ID")
	}

	fullStructs = []testtypes.FullStruct{}

	if n, _ := db.Find(&fullStructs).Run(); n != 2 {
		t.Errorf("Testing querying with 2 entity saved. Expecting 2 entities - got %v", n)
	}

	if c, _ := db.Find(&fullStructs).Count(); c != 2 {
		t.Errorf("Testing count 2 entities saved. Expecting 2 - got %v", c)
	}

	if fullStructs[0].ID == fullStructs[1].ID {
		t.Errorf("Testing querying with 2 entities saved. 2 results returned. Both have same ID")
	}

	// Limit
	fullStructs = []testtypes.FullStruct{}
	if n, _ := db.Find(&fullStructs).Limit(1).Run(); n != 1 {
		t.Errorf("Testing querying with 2 entities saved + limit. Wrong number of results received")
	}

	// Reverse - simple, only tests number received
	fullStructs = []testtypes.FullStruct{}
	if n, _ := db.Find(&fullStructs).Reverse().Run(); n != 2 {
		t.Errorf("Testing querying with 2 entities saved + reverse. Expected %v, got %v", 2, n)
	}

	// Reverse + Limit - simple, only tests number received
	fullStructs = []testtypes.FullStruct{}
	if n, _ := db.Find(&fullStructs).Reverse().Limit(1).Run(); n != 1 {
		t.Errorf("Testing querying with 2 entities saved + reverse + limit. Expected %v, got %v", 1, n)
	}

	// Reverse + Count
	fullStructs = []testtypes.FullStruct{}
	if c, _ := db.Find(&fullStructs).Reverse().Count(); c != 2 {
		t.Errorf("Testing count with 2 entities saved + reverse. Expected %v, got %v", 2, c)
	}

	// Compare forwards and backwards
	forwards := []testtypes.FullStruct{}
	backwards := []testtypes.FullStruct{}
	db.Find(&forwards).Run()
	db.Find(&backwards).Reverse().Run()
	if forwards[0].ID != backwards[1].ID || forwards[1].ID != backwards[0].ID {
		t.Error("Comparing regular and reversed results. Fist and last of each list should be the same but were not")
	}

}

func Test_BasicQuery_First(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	tt1 := testtypes.FullStruct{}
	tt2 := testtypes.FullStruct{}
	db.Save(&tt1, &tt2)

	var fullStruct testtypes.FullStruct
	n, err := db.First(&fullStruct).Run()

	if err != nil {
		t.Error("Testing first - got error")
	}

	if n != 1 {
		t.Errorf("Testing first. Expecting 1 entity - got %v", n)
	}

	if fullStruct.ID.IsNil() {
		t.Errorf("Testing first. Nil ID retrieved")
	}

	if fullStruct.ID != tt1.ID {
		t.Errorf("Testing first. Order IDs are not equal - wrong fullStruct retrieved")
	}

	// Test nothing found (impossible range)
	n, _ = db.First(&fullStruct).From(time.Now()).To(time.Now()).Run()
	if n != 0 {
		t.Errorf("Testing first when nothing should be found.  Got n = %v", n)
	}
}
