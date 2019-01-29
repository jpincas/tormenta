package tormenta

import (
	"os"
	"testing"
)

func Test_Open_ValidDirectory(t *testing.T) {
	testName := "Testing opening Torment DB connection with a valid directory"
	dir := "data/test"

	// Create a connection to a test DB
	db, err := Open(dir, DefaultOptions)
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
	testName := "Testing closing TormentaDB connection"

	db, _ := Open("data/test", DefaultOptions)
	err := db.Close()
	if err != nil {
		t.Errorf("%s. Failed to close connection with error: %v", testName, err)
	}
}

func Test_Open_InvalidDirectory(t *testing.T) {
	testName := "Testing opening Torment DB connection with an invalid directory"

	// Create a connection to a test DB
	db, err := Open("", DefaultOptions)

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
	db, err := OpenTest(dir, DefaultOptions)

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
