package tormenta

import (
	"errors"
	"os"

	"github.com/dgraph-io/badger"
)

// DB is the wrapper of the BadgerDB connection
type DB struct {
	KV *badger.DB
}

// testDirectory alters a specified data directory to mark it as for tests
func testDirectory(dir string) string {
	return dir + "-test"
}

// Open returns a connection to TormentDB connection
func Open(dir string) (*DB, error) {
	if dir == "" {
		return nil, errors.New("No valid data directory provided")
	}

	// Create directory if does not exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, errors.New("Could not create data directory")
		}
	}

	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = dir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// OpenTest is a convenience function to wipe the existing data at the specified location and create a new connection.  As a safety measure against production use, it will append "-test" to the directory name
func OpenTest(dir string) (*DB, error) {
	testDir := testDirectory(dir)

	// Attempt to remove the existing directory
	os.RemoveAll("./" + testDir)

	// Now check if it exists
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		return nil, errors.New("Could not remove existing data directory")
	}

	return Open(testDir)
}

// Close closes the connection to the DB
func (db DB) Close() error {
	return db.KV.Close()
}
