package tormenta_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/testtypes"
)

var zeroValueTime time.Time

func Test_BasicSave(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	// Create basic testtypes.FullStruct and save
	fullStruct := testtypes.FullStruct{}
	n, err := db.Save(&fullStruct)

	// Test any error
	if err != nil {
		t.Errorf("Testing basic record save. Got error: %v", err)
	}

	// Test that 1 record was reported saved
	if n != 1 {
		t.Errorf("Testing basic record save. Expected 1 record saved, got %v", n)
	}

	// Check ID has been set
	if fullStruct.ID.IsNil() {
		t.Error("Testing basic record save with create new ID. ID after save is nil")
	}

	//  Check that updated field was set
	if fullStruct.LastUpdated == zeroValueTime {
		t.Error("Testing basic record save. 'Last Upated' is time zero value")
	}

	// Take a snapshot
	fullStructBeforeSecondSave := fullStruct

	// Save again
	n2, err2 := db.Save(&fullStruct)

	// Basic tests
	if err2 != nil {
		t.Errorf("Testing 2nd record save. Got error %v", err)
	}

	if n2 != 1 {
		t.Errorf("Testing 2nd record save. Expected 1 record saved, got %v", n)
	}

	//  Check that updated field was updated:the new value
	// should obviously be later
	if !fullStructBeforeSecondSave.LastUpdated.Before(fullStruct.LastUpdated) {
		t.Error("Testing 2nd record save. 'Created' time has changed")
	}
}

func Test_SaveDifferentTypes(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	// Create basic testtypes.FullStruct and save
	fullStruct := testtypes.FullStruct{}
	miniStruct := testtypes.MiniStruct{}
	n, err := db.Save(&fullStruct, &miniStruct)

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
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	// Create basic testtypes.FullStruct and save
	fullStruct := testtypes.FullStruct{}
	db.Save(&fullStruct)

	// Test postsave trigger
	if !fullStruct.IsSaved {
		t.Error("Testing postsave trigger.  isSaved should be true but was not")
	}

	// Set up a condition that will cause the testtypes.FullStruct not to save
	fullStruct.ShouldBlockSave = true

	// Test presave trigger
	n, err := db.Save(&fullStruct)
	if fullStruct.TriggerString != "triggered" {
		t.Errorf("Testing presave trigger.  TriggerStringField wrong. Expected %s, got %s", "triggered", fullStruct.TriggerString)
	}

	if n != 0 || err == nil {
		t.Error("Testing presave trigger.  This record should not have saved, but it did and no error returned")
	}
}

type structA struct {
	StringField string
	tormenta.Model
}

type structB struct {
	StringField string
	tormenta.Model
}

func (s *structA) PreSave(db tormenta.DB) ([]tormenta.Record, error) {
	return []tormenta.Record{&structB{StringField: "b"}}, nil
}

func Test_SaveTrigger_CascadingSaves(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	// Saving A should also create and save a B according to the presave trigger
	if _, err := db.Save(&structA{}); err != nil {
		t.Errorf("Testing presave trigger with cascades. Got error %v", err)
	}

	// So lets see if its there
	res := structB{}
	if n, err := db.First(&res).Run(); err != nil {
		t.Errorf("Testing presave trigger with cascades. Got error %v", err)
	} else if n != 1 {
		t.Errorf("Testing presave trigger with cascades. Trying to retrieve struct B, but got n=%v", n)
	}

	if res.StringField != "b" {
		t.Errorf("Testing presave trigger with cascades. Checking struct B string value, expected %s but got %s", "b", res.StringField)
	}
}

func Test_SaveMultiple(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	fullStruct1 := testtypes.FullStruct{}
	fullStruct2 := testtypes.FullStruct{}

	// Multiple argument syntax
	n, _ := db.Save(&fullStruct1, &fullStruct2)
	if n != 2 {
		t.Errorf("Testing multiple save. Expected %v, got %v", 2, n)
	}

	if fullStruct1.ID == fullStruct2.ID {
		t.Errorf("Testing multiple save. 2 testtypes.FullStructs have same ID")
	}

	// Spread syntax
	// A little akward as you can't just pass in the slice of entities
	// You have to manually translate to []Record
	n, _ = db.Save([]tormenta.Record{&fullStruct1, &fullStruct2}...)
	if n != 2 {
		t.Errorf("Testing multiple save. Expected %v, got %v", 2, n)
	}

}

func Test_SaveMultipleLarge(t *testing.T) {
	const noOfTests = 500

	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	var fullStructsToSave []tormenta.Record

	for i := 0; i < noOfTests; i++ {
		fullStructsToSave = append(fullStructsToSave, &testtypes.FullStruct{
			StringField: fmt.Sprintf("customer-%v", i),
		})
	}

	n, err := db.Save(fullStructsToSave...)
	if err != nil {
		t.Errorf("Testing save large number of entities. Got error: %s", err)
	}

	if n != noOfTests {
		t.Errorf("Testing save large number of entities. Expected %v, got %v.  Err: %s", noOfTests, n, err)
	}

	var fullStructs []testtypes.FullStruct
	n, _ = db.Find(&fullStructs).Run()
	if n != noOfTests {
		t.Errorf("Testing save large number of entities, then retrieve. Expected %v, got %v", noOfTests, n)
	}

}

// Badger can only take a certain number of entities per transaction -
// which depends on how large the entities are.
// It should give back an error if we try to save too many
func Test_SaveMultipleTooLarge(t *testing.T) {
	const noOfTests = 1000000

	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	var fullStructsToSave []tormenta.Record

	for i := 0; i < noOfTests; i++ {
		fullStructsToSave = append(fullStructsToSave, &testtypes.FullStruct{})
	}

	n, err := db.Save(fullStructsToSave...)
	if err == nil {
		t.Error("Testing save large number of entities.Expecting an error but did not get one")

	}

	if n != 0 {
		t.Errorf("Testing save large number of entities. Expected %v, got %v", 0, n)
	}

	var fullStructs []testtypes.FullStruct
	n, _ = db.Find(&fullStructs).Run()
	if n != 0 {
		t.Errorf("Testing save large number of entities, then retrieve. Expected %v, got %v", 0, n)
	}

}

func Test_SaveMultipleLargeIndividually(t *testing.T) {
	const noOfTests = 10000

	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	var fullStructsToSave []tormenta.Record

	for i := 0; i < noOfTests; i++ {
		fullStructsToSave = append(fullStructsToSave, &testtypes.FullStruct{})
	}

	n, err := db.SaveIndividually(fullStructsToSave...)
	if err != nil {
		t.Errorf("Testing save large number of entities individually. Got error: %s", err)
	}

	if n != noOfTests {
		t.Errorf("Testing save large number of entities. Expected %v, got %v", 0, n)
	}
}

func Test_Save_SkipFields(t *testing.T) {
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
	defer db.Close()

	// Create basic testtypes.FullStruct and save
	fullStruct := testtypes.FullStruct{
		// Include a field that shouldnt be deleted
		IntField:                    1,
		NoSaveSimple:                "somthing",
		NoSaveTwoTags:               "somthing",
		NoSaveTwoTagsDifferentOrder: "somthing",
		NoSaveJSONSkiptag:           "something",

		// This one changes the name of the JSON tag
		NoSaveJSONtag: "somthing",
	}
	n, err := db.Save(&fullStruct)

	// Test any error
	if err != nil {
		t.Errorf("Testing save with skip field. Got error: %v", err)
	}

	// Test that 1 record was reported saved
	if n != 1 {
		t.Errorf("Testing save with skip field. Expected 1 record saved, got %v", n)
	}

	// Read back the record into a different target
	var readRecord testtypes.FullStruct
	found, err := db.Get(&readRecord, fullStruct.ID)

	// Test any error
	if err != nil {
		t.Errorf("Testing save with skip field. Got error reading back: %v", err)
	}

	// Test that 1 record was read back
	if !found {
		t.Errorf("Testing save with skip field. Expected 1 record read back, got %v", n)
	}

	// Test all the fields that should not have been saved
	if readRecord.IntField != 1 {
		t.Error("Testing save with skip field. Looks like IntField was deleted when it shouldnt have been")
	}

	if readRecord.NoSaveSimple != "" {
		t.Errorf("Testing save with skip field. NoSaveSimple should have been blank but was '%s'", readRecord.NoSaveSimple)
	}

	if readRecord.NoSaveTwoTags != "" {
		t.Errorf("Testing save with skip field. NoSaveTwoTags should have been blank but was '%s'", readRecord.NoSaveTwoTags)
	}

	if readRecord.NoSaveTwoTagsDifferentOrder != "" {
		t.Errorf("Testing save with skip field. NoSaveTwoTagsDifferentOrder should have been blank but was '%s'", readRecord.NoSaveTwoTagsDifferentOrder)
	}

	if readRecord.NoSaveJSONtag != "" {
		t.Errorf("Testing save with skip field. NoSaveJSONtag should have been blank but was '%s'", readRecord.NoSaveJSONtag)
	}

	if readRecord.NoSaveJSONSkiptag != "" {
		t.Errorf("Testing save with skip field. NoSaveJSONSkiptag should have been blank but was '%s'", readRecord.NoSaveJSONSkiptag)
	}
}
