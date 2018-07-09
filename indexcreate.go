package tormenta

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

// Index format
// i:indexname:root:indexcontent:entityID
// i:order:customer:5:324ds-3werwf-234wef-23wef

const (
	tormentaTag      = "tormenta"
	tormentaTagIndex = "index"
)

func index(txn *badger.Txn, entity Tormentable, keyRoot []byte, id gouuidv6.UUID) error {
	v := reflect.Indirect(reflect.ValueOf(entity))

	for i := 0; i < v.NumField(); i++ {

		// Look for the 'tormenta:index' tag
		fieldType := v.Type().Field(i)
		if idx := fieldType.Tag.Get(tormentaTag); idx == tormentaTagIndex {

			// Create the index key
			key := makeIndexKey(
				keyRoot,
				id,
				fieldType.Name,
				v.Field(i).Interface(),
			)

			// Set the index key with no content
			if err := txn.Set(key, []byte{}); err != nil {
				return err
			}
		}
	}

	return nil
}

func makeIndexKey(root []byte, id gouuidv6.UUID, indexName string, indexContent interface{}) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(indexKeyPrefix),
			root,
			[]byte(strings.ToLower(indexName)),
			interfaceToBytes(indexContent),
			id.Bytes(),
		},
		[]byte(keySeparator),
	)
}

func interfaceToBytes(value interface{}) []byte {
	// Must use BigEndian for correct sorting

	// Empty interface -> empty byte slice
	if value == nil {
		return []byte{}
	}

	// Set up buffer for writing binary values
	buf := new(bytes.Buffer)

	switch value.(type) {
	// For ints, cast the interface to int, convert to uint32 (normal ints don't work)
	case int:
		binary.Write(buf, binary.BigEndian, uint32(value.(int)))
		return buf.Bytes()
	// For floats, write straight to binary
	case float64:
		binary.Write(buf, binary.BigEndian, value.(float64))
		return buf.Bytes()
	}

	// Everything else as a string
	return []byte(fmt.Sprintf("%v", value))
}
