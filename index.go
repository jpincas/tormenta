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

func index(txn *badger.Txn, entity Tormentable, keyRoot []byte, id gouuidv6.UUID) error {
	v := reflect.Indirect(reflect.ValueOf(entity))

	for i := 0; i < v.NumField(); i++ {

		// Look for the 'tormenta:index' tag
		fieldType := v.Type().Field(i)
		if idx := fieldType.Tag.Get("tormenta"); idx == "index" {

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
			[]byte(strings.ToLower(indexName)),
			root,
			interfaceToBytes(indexContent),
			id.Bytes(),
		},
		[]byte(":"),
	)
}

func interfaceToBytes(value interface{}) []byte {
	buf := new(bytes.Buffer)
	var toWrite interface{}

	switch value.(type) {
	case int:
		toWrite = value.(int)
	case string:
		toWrite = value.(string)
	case float64:
		toWrite = value.(float64)
	default:
		return []byte(fmt.Sprintf("%v", value))
	}

	binary.Write(buf, binary.LittleEndian, toWrite)
	return []byte(buf.Bytes())
}
