package archive_test

// TODO: bug in the delelter with std lib json

// func Test_Save_SkipFields(t *testing.T) {
// 	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
// 	defer db.Close()

// 	// Create basic testtypes.FullStruct and save
// 	fullStruct := testtypes.FullStruct{
// 		// Include a field that shouldnt be deleted
// 		IntField:                    1,
// 		NoSaveSimple:                "somthing",
// 		NoSaveTwoTags:               "somthing",
// 		NoSaveTwoTagsDifferentOrder: "somthing",
// 		NoSaveJSONSkiptag:           "something",

// 		// This one changes the name of the JSON tag
// 		NoSaveJSONtag: "somthing",
// 	}
// 	n, err := db.Save(&fullStruct)

// 	// Test any error
// 	if err != nil {
// 		t.Errorf("Testing save with skip field. Got error: %v", err)
// 	}

// 	// Test that 1 record was reported saved
// 	if n != 1 {
// 		t.Errorf("Testing save with skip field. Expected 1 record saved, got %v", n)
// 	}

// 	// Read back the record into a different target
// 	var readRecord testtypes.FullStruct
// 	found, err := db.Get(&readRecord, fullStruct.ID)

// 	// Test any error
// 	if err != nil {
// 		t.Errorf("Testing save with skip field. Got error reading back: %v", err)
// 	}

// 	// Test that 1 record was read back
// 	if !found {
// 		t.Errorf("Testing save with skip field. Expected 1 record read back, got %v", n)
// 	}

// 	// Test all the fields that should not have been saved
// 	if readRecord.IntField != 1 {
// 		t.Error("Testing save with skip field. Looks like IntField was deleted when it shouldnt have been")
// 	}

// 	if readRecord.NoSaveSimple != "" {
// 		t.Errorf("Testing save with skip field. NoSaveSimple should have been blank but was '%s'", readRecord.NoSaveSimple)
// 	}

// 	if readRecord.NoSaveTwoTags != "" {
// 		t.Errorf("Testing save with skip field. NoSaveTwoTags should have been blank but was '%s'", readRecord.NoSaveTwoTags)
// 	}

// 	if readRecord.NoSaveTwoTagsDifferentOrder != "" {
// 		t.Errorf("Testing save with skip field. NoSaveTwoTagsDifferentOrder should have been blank but was '%s'", readRecord.NoSaveTwoTagsDifferentOrder)
// 	}

// 	if readRecord.NoSaveJSONtag != "" {
// 		t.Errorf("Testing save with skip field. NoSaveJSONtag should have been blank but was '%s'", readRecord.NoSaveJSONtag)
// 	}

// 	if readRecord.NoSaveJSONSkiptag != "" {
// 		t.Errorf("Testing save with skip field. NoSaveJSONSkiptag should have been blank but was '%s'", readRecord.NoSaveJSONSkiptag)
// 	}
// }
