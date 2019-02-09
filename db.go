package tormenta

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
	jsoniter "github.com/json-iterator/go"
)

const (
	// Serialisers
	SerialiserJSONStdLib = "json"
	SerialiserJSONIter   = "jsoniter"
)

// DB is the wrapper of the BadgerDB connection
type DB struct {
	KV              *badger.DB
	serialiseFunc   func(interface{}) ([]byte, error)
	unserialiseFunc func([]byte, interface{}) error
}

type Options struct {
	JsonIterAPI jsoniter.API
	Serialiser  string
}

var DefaultOptions = Options{
	// Use the fastest JSONiter option by default
	// Main difference is precision of floats - see https://godoc.org/github.com/json-iterator/go
	JsonIterAPI: jsoniter.ConfigFastest,
	Serialiser:  SerialiserJSONIter,
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
	badgerDB, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return openDB(badgerDB, options)
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

func openDB(badgerDB *badger.DB, options Options) (*DB, error) {
	db := &DB{
		KV: badgerDB,
	}

	switch options.Serialiser {
	case SerialiserJSONIter:
		db.serialiseFunc = options.JsonIterAPI.Marshal
		db.unserialiseFunc = options.JsonIterAPI.Unmarshal
	case SerialiserJSONStdLib:
		db.serialiseFunc = json.Marshal
		db.unserialiseFunc = json.Unmarshal

	default:
		return db, fmt.Errorf("%s is not a valid serialiser", options.Serialiser)
	}

	return db, nil
}

func (db DB) unserialise(val []byte, entity Record) error {
	return db.unserialiseFunc(val, entity)
}

func (db DB) serialise(entity Record) ([]byte, error) {
	return db.serialiseFunc(entity)
}
