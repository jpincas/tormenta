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
// i:fullStruct:customer:5:324ds-3werwf-234wef-23wef

const (
	tormentaTag        = "tormenta"
	tormentaTagNoIndex = "noindex"
	tormentaTagSplit   = "split"
)

var (
	typeInt    = reflect.TypeOf(0)
	typeFloat  = reflect.TypeOf(0.99)
	typeString = reflect.TypeOf("")
)

func index(txn *badger.Txn, entity Record) error {
	keys := indexStruct(
		recordValue(entity),
		entity,
		KeyRoot(entity),
		entity.GetID(),
	)

	for i := range keys {
		if err := txn.Set(keys[i], []byte{}); err != nil {
			return err
		}
	}

	return nil
}

func indexStruct(v reflect.Value, entity Record, keyRoot []byte, id gouuidv6.UUID) (keys [][]byte) {
	for i := 0; i < v.NumField(); i++ {

		// If the 'tormenta:noindex' tag is present, don't index
		fieldType := v.Type().Field(i)
		if idx := fieldType.Tag.Get(tormentaTag); idx != tormentaTagNoIndex {
			switch fieldType.Type.Kind() {
			case reflect.Slice:
				keys = append(keys, getMultipleIndexKeys(v.Field(i), keyRoot, id, fieldType.Name)...)

			case reflect.Array:
				keys = append(keys, getMultipleIndexKeys(v.Field(i), keyRoot, id, fieldType.Name)...)

			case reflect.String:
				// If the string is tagged with 'split',
				// then index each of the words separately
				if idx == tormentaTagSplit {
					keys = append(keys, getSplitStringIndexes(v.Field(i), keyRoot, id, fieldType.Name)...)
				} else {
					keys = append(keys, makeIndexKey(keyRoot, id, fieldType.Name, v.Field(i).Interface()))
				}

			case reflect.Struct:
				// Recursively index embedded structs
				if fieldType.Anonymous {
					keys = append(keys, indexStruct(v.Field(i), entity, keyRoot, id)...)
				}

			default:
				keys = append(keys, makeIndexKey(keyRoot, id, fieldType.Name, v.Field(i).Interface()))
			}
		}
	}

	return
}

// MakeIndexKey constructs an index key
func MakeIndexKey(root []byte, id gouuidv6.UUID, indexName string, indexContent interface{}) []byte {
	return makeIndexKey(root, id, indexName, indexContent)
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

func getMultipleIndexKeys(v reflect.Value, root []byte, id gouuidv6.UUID, indexName string) (keys [][]byte) {
	for i := 0; i < v.Len(); i++ {
		key := makeIndexKey(root, id, indexName, v.Index(i).Interface())
		keys = append(keys, key)
	}

	return
}

func getSplitStringIndexes(v reflect.Value, root []byte, id gouuidv6.UUID, indexName string) (keys [][]byte) {
	strings := strings.Split(v.String(), " ")

	// Clean non-content words
	strings = removeNonContentWords(strings)

	for _, s := range strings {
		key := makeIndexKey(root, id, indexName, s)
		keys = append(keys, key)
	}

	return
}

func interfaceToBytes(value interface{}) []byte {
	// Note: must use BigEndian for correct sorting

	// Empty interface -> empty byte slice
	if value == nil {
		return []byte{}
	}

	// Set up buffer for writing binary values
	buf := new(bytes.Buffer)

	// Rather than doing a simple type switch .(type),
	// we switch on the Kind
	// This provides a neat solution to 'defined' or 'named' types
	// e.g. the kind of 'string' is 'string'
	// AND the kind of NamedString (with underlying type string) is ALSO just string

	switch reflect.ValueOf(value).Type().Kind() {

	// For ints, cast the interface to int, convert to uint32 (normal ints don't work)
	case reflect.Int:
		i := convertUnderlying(value, typeInt)
		binary.Write(buf, binary.BigEndian, uint32(i.(int)))

		return buf.Bytes()

	// For floats, write straight to binary
	case reflect.Float64:
		i := convertUnderlying(value, typeFloat)
		binary.Write(buf, binary.BigEndian, i.(float64))

		return buf.Bytes()

	// For strings, lower case before indexing
	case reflect.String:
		i := convertUnderlying(value, typeString)
		return []byte(strings.ToLower(i.(string)))
	}

	// Everything else as a string
	return []byte(fmt.Sprintf("%v", value))
}

// convertUnderlying takes an interface and converts its underlying type
// to the target type.  Obviously the underlying types must be convertible
// E.g. NamedInt -> Int
func convertUnderlying(src interface{}, targetType reflect.Type) interface{} {
	return reflect.ValueOf(src).Convert(targetType).Interface()
}
