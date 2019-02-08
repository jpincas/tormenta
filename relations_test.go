package tormenta_test

import (
	"testing"

	"github.com/jpincas/tormenta/testtypes"

	"github.com/jpincas/tormenta"
)

func Test_Relations_HasOne(t *testing.T) {
	// Open the DB
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
	defer db.Close()

	doubleNestedStruct1 := testtypes.DoubleNestedRelatedStruct{}
	doubleNestedStruct2 := testtypes.DoubleNestedRelatedStruct{}
	db.Save(&doubleNestedStruct1, &doubleNestedStruct2)

	// Created some nested structs and save
	nestedStruct1 := testtypes.NestedRelatedStruct{
		NestedID: doubleNestedStruct2.ID,
	}
	nestedStruct2 := testtypes.NestedRelatedStruct{
		NestedID: doubleNestedStruct1.ID,
	}
	db.Save(&nestedStruct1, &nestedStruct2)

	// Create some related structs which nest the above and save
	relatedStruct1 := testtypes.RelatedStruct{
		NestedID: nestedStruct2.ID,
	}
	relatedStruct2 := testtypes.RelatedStruct{
		NestedID: nestedStruct1.ID,
	}
	db.Save(&relatedStruct1, &relatedStruct2)

	// Create some full structs including these relations
	struct1 := testtypes.FullStruct{
		HasOneID:        relatedStruct1.ID,
		HasAnotherOneID: relatedStruct1.ID,
	}

	struct2 := testtypes.FullStruct{
		HasOneID:        relatedStruct2.ID,
		HasAnotherOneID: relatedStruct2.ID,
	}

	struct3 := testtypes.FullStruct{
		HasOneID:        relatedStruct1.ID,
		HasAnotherOneID: relatedStruct1.ID,
	}
	db.Save(&struct1, &struct2, &struct3)

	// A) Single invalid relation, no nesting
	// Reload
	fullStructs := []testtypes.FullStruct{}
	if n, err := db.Find(&fullStructs).Run(); err != nil || n != 3 {
		t.Errorf("Save/retrieve failed. Err: %v; n: %v", err, n)
	}

	// Convert to records
	records := []tormenta.Record{}
	for i := range fullStructs {
		records = append(records, &fullStructs[i])
	}

	testName := "non existent relation"
	if err := tormenta.HasOne(db, []string{"DoesntHaveOne"}, records...); err == nil {
		t.Errorf("Testing %s. Should have error but didn't", testName)
	}

	// 1) Single relation, no nesting
	// Reload
	fullStructs = []testtypes.FullStruct{}
	if n, err := db.Find(&fullStructs).Run(); err != nil || n != 3 {
		t.Errorf("Save/retrieve failed. Err: %v; n: %v", err, n)
	}

	// Convert to records
	records = []tormenta.Record{}
	for i := range fullStructs {
		records = append(records, &fullStructs[i])
	}

	testName = "single, one level relation"
	if err := tormenta.HasOne(db, []string{"HasOne"}, records...); err != nil {
		t.Errorf("Testing %s. Error loading relations: %s", testName, err)
	}

	for i, fullStruct := range fullStructs {
		if fullStruct.HasOneID != fullStruct.HasOne.ID {
			t.Errorf(
				"Testing %s. Comparing HasOneID to HasOne.ID for order %v and they are not the same: %v vs %v",
				testName,
				i,
				fullStruct.HasOneID,
				fullStruct.HasOne.ID,
			)
		}
	}

	// 2) Two relations, but repeated
	// Reload
	fullStructs = []testtypes.FullStruct{}
	if n, err := db.Find(&fullStructs).Run(); err != nil || n != 3 {
		t.Errorf("Save/retrieve failed. Err: %v; n: %v", err, n)
	}

	// Convert to records
	records = []tormenta.Record{}
	for i := range fullStructs {
		records = append(records, &fullStructs[i])
	}

	testName = "two relations, one level, repeated"
	if err := tormenta.HasOne(db, []string{"HasOne", "HasOne"}, records...); err != nil {
		t.Errorf("Testing %s. Error loading relations: %s", testName, err)
	}

	for i, fullStruct := range fullStructs {
		if fullStruct.HasOneID != fullStruct.HasOne.ID {
			t.Errorf(
				"Testing %s. Comparing HasOneID to HasOne.ID for order %v and they are not the same: %v vs %v",
				testName,
				i,
				fullStruct.HasOneID,
				fullStruct.HasOne.ID,
			)
		}
	}

	// 3) Two relations
	// Reload
	fullStructs = []testtypes.FullStruct{}
	if n, err := db.Find(&fullStructs).Run(); err != nil || n != 3 {
		t.Errorf("Save/retrieve failed. Err: %v; n: %v", err, n)
	}

	// Convert to records
	records = []tormenta.Record{}
	for i := range fullStructs {
		records = append(records, &fullStructs[i])
	}

	testName = "two relations, one level, repeated"
	if err := tormenta.HasOne(db, []string{"HasOne", "HasAnotherOne"}, records...); err != nil {
		t.Errorf("Testing %s. Error loading relations: %s", testName, err)
	}

	for i, fullStruct := range fullStructs {
		if fullStruct.HasOneID != fullStruct.HasOne.ID {
			t.Errorf(
				"Testing %s. Comparing HasOneID to HasOne.ID for order %v and they are not the same: %v vs %v",
				testName,
				i,
				fullStruct.HasOneID,
				fullStruct.HasOne.ID,
			)
		}

		if fullStruct.HasAnotherOneID != fullStruct.HasAnotherOne.ID {
			t.Errorf(
				"Testing %s. Comparing HasAnotherOneID to HasAnotherOne.ID for order %v and they are not the same: %v vs %v",
				testName,
				i,
				fullStruct.HasAnotherOneID,
				fullStruct.HasAnotherOne.ID,
			)
		}
	}

	// 4) Single relation, 2 level nesting
	// Reload
	fullStructs = []testtypes.FullStruct{}
	if n, err := db.Find(&fullStructs).Run(); err != nil || n != 3 {
		t.Errorf("Save/retrieve failed. Err: %v; n: %v", err, n)
	}

	// Convert to records
	records = []tormenta.Record{}
	for i := range fullStructs {
		records = append(records, &fullStructs[i])
	}

	testName = "single, nested relation"
	if err := tormenta.HasOne(db, []string{"HasOne.Nested"}, records...); err != nil {
		t.Errorf("Testing %s. Error loading relations: %s", testName, err)
	}

	for i, fullStruct := range fullStructs {
		if fullStruct.HasOneID != fullStruct.HasOne.ID {
			t.Errorf(
				"Testing %s. Comparing HasOneID to HasOne.ID for order %v and they are not the same: %v vs %v",
				testName,
				i,
				fullStruct.HasOneID,
				fullStruct.HasOne.ID,
			)
		}

		if fullStruct.HasOne.Nested == nil {
			t.Fatalf("Testing %s - index %v. HasOne.Nested is nil - the relation didn't load", testName, i)
		}

		if fullStruct.HasOne.NestedID != fullStruct.HasOne.Nested.ID {
			t.Errorf(
				"Testing %s. Comparing HasOne.NestedID to HasOne.Nested.ID for order %v and they are not the same: %v vs %v",
				testName,
				i,
				fullStruct.HasOne.NestedID,
				fullStruct.HasOne.Nested.ID,
			)
		}
	}

	// 4) Single relation, 3 level nesting
	// Reload
	fullStructs = []testtypes.FullStruct{}
	if n, err := db.Find(&fullStructs).Run(); err != nil || n != 3 {
		t.Errorf("Save/retrieve failed. Err: %v; n: %v", err, n)
	}

	// Convert to records
	records = []tormenta.Record{}
	for i := range fullStructs {
		records = append(records, &fullStructs[i])
	}

	testName = "single, 3 level nested relation"
	if err := tormenta.HasOne(db, []string{"HasOne.Nested.Nested"}, records...); err != nil {
		t.Errorf("Testing %s. Error loading relations: %s", testName, err)
	}

	for i, fullStruct := range fullStructs {
		if fullStruct.HasOneID != fullStruct.HasOne.ID {
			t.Errorf(
				"Testing %s. Comparing HasOneID to HasOne.ID for order %v and they are not the same: %v vs %v",
				testName,
				i,
				fullStruct.HasOneID,
				fullStruct.HasOne.ID,
			)
		}

		// L1

		if fullStruct.HasOne.Nested == nil {
			t.Fatalf("Testing %s - index %v. HasOne.Nested is nil - the relation didn't load", testName, i)
		}

		if fullStruct.HasOne.NestedID != fullStruct.HasOne.Nested.ID {
			t.Errorf(
				"Testing %s. Comparing HasOne.NestedID to HasOne.Nested.ID for order %v and they are not the same: %v vs %v",
				testName,
				i,
				fullStruct.HasOne.NestedID,
				fullStruct.HasOne.Nested.ID,
			)
		}

		// L2

		if fullStruct.HasOne.Nested.Nested == nil {
			t.Fatalf("Testing %s - index %v. HasOne.Nested.Nested is nil - the relation didn't load", testName, i)
		}

		if fullStruct.HasOne.Nested.NestedID != fullStruct.HasOne.Nested.Nested.ID {
			t.Errorf(
				"Testing %s. Comparing HasOne.Nested.NestedID to HasOne.Nested.Nested.ID for order %v and they are not the same: %v vs %v",
				testName,
				i,
				fullStruct.HasOne.NestedID,
				fullStruct.HasOne.Nested.ID,
			)
		}

	}
}
