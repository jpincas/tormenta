package tormenta_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jpincas/tormenta"
)

var zeroValueTime time.Time

func Test_BasicSave(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// Create basic testType and save
	testType := TestType{}
	n, err := db.Save(&testType)

	// Test any error
	if err != nil {
		t.Errorf("Testing basic record save. Got error: %v", err)
	}

	// Test that 1 record was reported saved
	if n != 1 {
		t.Errorf("Testing basic record save. Expected 1 record saved, got %v", n)
	}

	// Check ID has been set
	if testType.ID.IsNil() {
		t.Error("Testing basic record save with create new ID. ID after save is nil")
	}

	//  Check that updated field was set
	if testType.LastUpdated == zeroValueTime {
		t.Error("Testing basic record save. 'Last Upated' is time zero value")
	}

	// Take a snapshot
	testTypeBeforeSecondSave := testType

	// Save again
	n2, err2 := db.Save(&testType)

	// Basic tests
	if err2 != nil {
		t.Errorf("Testing 2nd record save. Got error %v", err)
	}

	if n2 != 1 {
		t.Errorf("Testing 2nd record save. Expected 1 record saved, got %v", n)
	}

	//  Check that updated field was updated:the new value
	// should obviously be later
	if !testTypeBeforeSecondSave.LastUpdated.Before(testType.LastUpdated) {
		t.Error("Testing 2nd record save. 'Created' time has changed")
	}
}

func Test_SaveDifferentTypes(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// Create basic testType and save
	testType := TestType{}
	testType2 := TestType2{}
	n, err := db.Save(&testType, &testType2)

	// Test any error
	if err != nil {
		t.Errorf("Testing different records save. Got error: %v", err)
	}

	// Test that 2 records was reported saved
	if n != 2 {
		t.Errorf("Testing basic record save. Expected 1 record saved, got %v", n)
	}
}

func Test_SaveTrigger(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	// Create basic testType and save
	testType := TestType{}
	db.Save(&testType)

	// Test postsave trigger
	if !testType.IsSaved {
		t.Error("Testing postsave trigger.  isSaved should be true but was not")
	}

	// Set up a condition that will cause the testType not to save
	testType.ShouldBlockSave = true

	// Test presave trigger
	n, err := db.Save(&testType)
	if n != 0 || err == nil {
		t.Error("Testing presave trigger.  This record should not have saved, but it did and no error returned")
	}

}

func Test_SaveMultiple(t *testing.T) {
	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	testType1 := TestType{}
	testType2 := TestType{}

	// Multiple argument syntax
	n, _ := db.Save(&testType1, &testType2)
	if n != 2 {
		t.Errorf("Testing multiple save. Expected %v, got %v", 2, n)
	}

	if testType1.ID == testType2.ID {
		t.Errorf("Testing multiple save. 2 testTypes have same ID")
	}

	// Spread syntax
	// A little akward as you can't just pass in the slice of entities
	// You have to manually translate to []Record
	var testTypesToSave []tormenta.Record
	testTypes := []TestType{testType1, testType2}

	for _, testType := range testTypes {
		testTypesToSave = append(testTypesToSave, &testType)
	}

	n, _ = db.Save(testTypesToSave...)
	if n != 2 {
		t.Errorf("Testing multiple save. Expected %v, got %v", 2, n)
	}

}

func Test_SaveMultipleLarge(t *testing.T) {
	const notestTypes = 1000

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var testTypesToSave []tormenta.Record

	for i := 0; i < notestTypes; i++ {
		testTypesToSave = append(testTypesToSave, &TestType{
			StringField: fmt.Sprintf("customer-%v", i),
		})
	}

	n, err := db.Save(testTypesToSave...)
	if err != nil {
		t.Errorf("Testing save large number of entities. Got error: %s", err)
	}

	if n != notestTypes {
		t.Errorf("Testing save large number of entities. Expected %v, got %v.  Err: %s", notestTypes, n, err)
	}

	var testTypes []TestType
	n, _, _ = db.Find(&testTypes).Run()
	if n != notestTypes {
		t.Errorf("Testing save large number of entities, then retrieve. Expected %v, got %v", notestTypes, n)
	}

}

// Badger can only take a certain number of entities per transaction -
// which depends on how large the entities are.
// It should give back an error if we try to save too many
func Test_SaveMultipleTooLarge(t *testing.T) {
	const notestTypes = 1000000

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var testTypesToSave []tormenta.Record

	for i := 0; i < notestTypes; i++ {
		testTypesToSave = append(testTypesToSave, &TestType{})
	}

	n, err := db.Save(testTypesToSave...)
	if err == nil {
		t.Error("Testing save large number of entities.Expecting an error but did not get one")

	}

	if n != 0 {
		t.Errorf("Testing save large number of entities. Expected %v, got %v", 0, n)
	}

	var testTypes []TestType
	n, _, _ = db.Find(&testTypes).Run()
	if n != 0 {
		t.Errorf("Testing save large number of entities, then retrieve. Expected %v, got %v", 0, n)
	}

}

func Test_SaveMultipleLargeIndividually(t *testing.T) {
	const notestTypes = 10000

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()

	var testTypesToSave []tormenta.Record

	for i := 0; i < notestTypes; i++ {
		testTypesToSave = append(testTypesToSave, &TestType{})
	}

	n, err := db.SaveIndividually(testTypesToSave...)
	if err != nil {
		t.Errorf("Testing save large number of entities individually. Got error: %s", err)
	}

	if n != notestTypes {
		t.Errorf("Testing save large number of entities. Expected %v, got %v", 0, n)
	}
}
