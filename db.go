package tormenta

import (
	"errors"
	"os"

	"github.com/dgraph-io/badger"
	jsoniter "github.com/json-iterator/go"
)

// DB is the wrapper of the BadgerDB connection
type DB struct {
	KV   *badger.DB
	json jsoniter.API
}

type Options struct {
	json jsoniter.API
}

var DefaultOptions = Options{
	// Use the fasted JSONiter option by default
	// Main difference is precision of floats - see https://godoc.org/github.com/json-iterator/go
	json: jsoniter.ConfigFastest,
}

// testDirectory alters a specified data directory to mark it as for tests
func testDirectory(dir string) string {
	return dir + "-test"
}

// Open returns a connection to TormentDB connection
func Open(dir string, options Options) (*DB, error) {
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

	return &DB{
		KV:   db,
		json: options.json,
	}, nil
}

// OpenTest is a convenience function to wipe the existing data at the specified location and create a new connection.  As a safety measure against production use, it will append "-test" to the directory name
func OpenTest(dir string, options Options) (*DB, error) {
	testDir := testDirectory(dir)

	// Attempt to remove the existing directory
	os.RemoveAll("./" + testDir)

	// Now check if it exists
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		return nil, errors.New("Could not remove existing data directory")
	}

	return Open(testDir, options)
}

// Close closes the connection to the DB
func (db DB) Close() error {
	return db.KV.Close()
}
