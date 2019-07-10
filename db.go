package tormenta

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/dgraph-io/badger"
)

// DB is the wrapper of the BadgerDB connection
type DB struct {
	KV       *badger.DB
	Options  Options
	TestMode bool
	Dir      string
}

type Options struct {
	SerialiseFunc   func(interface{}) ([]byte, error)
	UnserialiseFunc func([]byte, interface{}) error
	BadgerOptions   badger.Options
	DebugMode       bool
}

var DefaultOptions = Options{
	SerialiseFunc:   json.Marshal,
	UnserialiseFunc: json.Unmarshal,
	BadgerOptions:   badger.DefaultOptions("data"),
	DebugMode:       false,
}

// testDirectory alters a specified data directory to mark it as for tests
func testDirectory(dir string) string {
	return dir + "-test"
}

func Open(dir string) (*DB, error) {
	return openWithOptions(dir, DefaultOptions, false)
}

// Open returns a connection to TormentDB connection
func OpenWithOptions(dir string, options Options) (*DB, error) {
	return openWithOptions(dir, options, false)
}

func OpenTest(dir string) (*DB, error) {
	return openTestWithOptions(dir, DefaultOptions)
}

// OpenTest is a convenience function to wipe the existing data at the specified location and create a new connection.  As a safety measure against production use, it will append "-test" to the directory name
func OpenTestWithOptions(dir string, options Options) (*DB, error) {
	return openTestWithOptions(dir, options)
}

func openTestWithOptions(dir string, options Options) (*DB, error) {
	testDir := testDirectory(dir)

	if err := removeDir(testDir); err != nil {
		return nil, err
	}

	return openWithOptions(testDir, options, true)
}

func openWithOptions(dir string, options Options, testMode bool) (*DB, error) {
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

	opts := options.BadgerOptions
	opts.Dir = dir
	opts.ValueDir = dir
	badgerDB, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &DB{
		KV:       badgerDB,
		Options:  options,
		Dir:      dir,
		TestMode: testMode,
	}, nil
}

// Close closes the connection to the DB
func (db DB) Close() error {
	// Close the DB connection first as it has a lock on the directory
	if err := db.KV.Close(); err != nil {
		return err
	}

	if db.TestMode {
		if err := removeDir(db.Dir); err != nil {
			return err
		}
	}

	return nil
}

func (db DB) unserialise(val []byte, entity interface{}) error {
	return db.Options.UnserialiseFunc(val, entity)
}

func (db DB) serialise(entity interface{}) ([]byte, error) {
	return db.Options.SerialiseFunc(entity)
}

func removeDir(dir string) error {
	// Attempt to remove the existing directory
	os.RemoveAll("./" + dir)

	// Now check if it exists
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return errors.New("Could not remove existing data directory")
	}

	return nil
}
