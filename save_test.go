package tormenta

import (
	"testing"
	"time"
)

var zeroValueTime time.Time

func Test_BasicSave(t *testing.T) {
	db, _ := OpenTest("data/tests")
	defer db.Close()

	// First Save

	// Create basic order and save
	order := Order{}
	n, err := db.Save(&order)

	// Test any error
	if err != nil {
		t.Errorf("Testing basic record save. Got error %v", err)
	}

	// Test that 1 record was reported saved
	if n != 1 {
		t.Errorf("Testing basic record save. Expected 1 record saved, got %v", n)
	}

	// Check ID has been set
	if order.ID.IsNil() {
		t.Error("Testing basic record save with create new ID. ID after save is nil")
	}

	//  Check that updated field was set
	if order.LastUpdated == zeroValueTime {
		t.Error("Testing basic record save. 'Last Upated' is time zero value")
	}

	// Take a snapshot
	orderBeforeSecondSave := order

	// Save again
	n2, err2 := db.Save(&order)

	// Basic tests
	if err2 != nil {
		t.Errorf("Testing 2nd record save. Got error %v", err)
	}

	if n2 != 1 {
		t.Errorf("Testing 2nd record save. Expected 1 record saved, got %v", n)
	}

	//  Check that updated field was updated:the new value
	// should obviously be later
	if !orderBeforeSecondSave.LastUpdated.Before(order.LastUpdated) {
		t.Error("Testing 2nd record save. 'Created' time has changed")
	}

}
