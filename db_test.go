package tormenta

import (
	"os"
	"testing"
)

func Test_Open_ValidDirectory(t *testing.T) {
	testName := "Testing opening Torment DB connection with a valid directory"
	dir := "data/testing"

	// Create a connection to a test DB
	db, err := OpenWithOptions(dir, DefaultOptions)
	defer db.Close()

	if err != nil {
		t.Errorf("%s. Failed to open connection with error %v", testName, err)
	}

	if db == nil {
		t.Errorf("%s. Failed to open connection. DB is nil", testName)
	}

	// Check the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("%s. Failed to create Torment data directory", testName)
	}

}

func Test_Close(t *testing.T) {
	dir := "data/test"
	testName := "Testing closing TormentaDB connection"

	db, _ := OpenWithOptions(dir, DefaultOptions)
	err := db.Close()
	if err != nil {
		t.Errorf("%s. Failed to close connection with error: %v", testName, err)
	}

	// Check the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("%s. Data directory vanished after close", testName)
	}
}

func Test_Open_InvalidDirectory(t *testing.T) {
	testName := "Testing opening Torment DB connection with an invalid directory"

	// Create a connection to a test DB
	db, err := OpenWithOptions("", DefaultOptions)

	if err == nil {
		t.Errorf("%s. Should have returned an error but did not", testName)
	}

	if db != nil {
		t.Errorf("%s. Should have returned a nil connection but did not", testName)
	}
}

func Test_Open_Test(t *testing.T) {
	testName := "Testing opening Torment DB with a blank DB"
	dir := "data/test"

	// Create a connection to a test DB
	db, err := OpenTestWithOptions(dir, DefaultOptions)

	if err != nil {
		t.Errorf("%s. Failed to open connection with error %v", testName, err)
	}

	if db == nil {
		t.Errorf("%s. Failed to open connection. DB is nil", testName)
	}

	// Check the directory exists
	if _, err := os.Stat(testDirectory(dir)); os.IsNotExist(err) {
		t.Errorf("%s. Failed to create Torment data directory", testName)
	}
}

func Test_Close_Test(t *testing.T) {
	dir := "data/test"
	testName := "Testing closing TormentaDB connection in test mode"

	db, _ := OpenTestWithOptions(dir, DefaultOptions)

	// Check the directory exists
	if _, err := os.Stat(testDirectory(dir)); os.IsNotExist(err) {
		t.Errorf("%s. Data directory not present", testName)
	}

	err := db.Close()
	if err != nil {
		t.Errorf("%s. Failed to close connection with error: %v", testName, err)
	}

	// Check the directory has been deleted
	if _, err := os.Stat(testDirectory(dir)); !os.IsNotExist(err) {
		t.Errorf("%s. Data directory should have been deleted after test", testName)
	}
}
