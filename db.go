package tormenta

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/dgraph-io/badger"
)

// DB is the wrapper of the BadgerDB connection
type DB struct {
	KV              *badger.DB
	Options         Options
	serialiseFunc   func(interface{}) ([]byte, error)
	unserialiseFunc func([]byte, interface{}) error
}

type Options struct {
	SerialiseFunc   func(interface{}) ([]byte, error)
	UnserialiseFunc func([]byte, interface{}) error
}

var DefaultOptions = Options{
	SerialiseFunc:   json.Marshal,
	UnserialiseFunc: json.Unmarshal,
}

// testDirectory alters a specified data directory to mark it as for tests
func testDirectory(dir string) string {
	return dir + "-test"
}

func Open(dir string) (*DB, error) {
	return OpenTestWithOptions(dir, DefaultOptions)
}

// Open returns a connection to TormentDB connection
func OpenWithOptions(dir string, options Options) (*DB, error) {
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
	badgerDB, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return openDB(badgerDB, options)
}

func OpenTest(dir string) (*DB, error) {
	testDir := testDirectory(dir)

	// Attempt to remove the existing directory
	os.RemoveAll("./" + testDir)

	// Now check if it exists
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		return nil, errors.New("Could not remove existing data directory")
	}

	return OpenWithOptions(testDir, DefaultOptions)
}

// OpenTest is a convenience function to wipe the existing data at the specified location and create a new connection.  As a safety measure against production use, it will append "-test" to the directory name
func OpenTestWithOptions(dir string, options Options) (*DB, error) {
	testDir := testDirectory(dir)

	// Attempt to remove the existing directory
	os.RemoveAll("./" + testDir)

	// Now check if it exists
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		return nil, errors.New("Could not remove existing data directory")
	}

	return OpenWithOptions(testDir, options)
}

// Close closes the connection to the DB
func (db DB) Close() error {
	return db.KV.Close()
}

func openDB(badgerDB *badger.DB, options Options) (*DB, error) {
	return &DB{
		KV:      badgerDB,
		Options: options,
	}, nil
}

func (db DB) unserialise(val []byte, entity interface{}) error {
	return db.Options.UnserialiseFunc(val, entity)
}

func (db DB) serialise(entity interface{}) ([]byte, error) {
	return db.Options.SerialiseFunc(entity)
}
