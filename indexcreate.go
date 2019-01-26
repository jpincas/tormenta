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
	tormentaTag        = "tormenta"
	tormentaTagNoIndex = "noindex"
	tormentaTagSplit   = "split"
)

var (
	typeInt    = reflect.TypeOf(0)
	typeFloat  = reflect.TypeOf(0.99)
	typeString = reflect.TypeOf("")
)

func index(txn *badger.Txn, entity Tormentable, keyRoot []byte, id gouuidv6.UUID) error {
	v := reflect.Indirect(reflect.ValueOf(entity))
	return indexStruct(v, txn, entity, keyRoot, id)
}

func indexStruct(v reflect.Value, txn *badger.Txn, entity Tormentable, keyRoot []byte, id gouuidv6.UUID) error {
	for i := 0; i < v.NumField(); i++ {

		// If the 'tormenta:noindex' tag is present, don't index
		fieldType := v.Type().Field(i)
		if idx := fieldType.Tag.Get(tormentaTag); idx != tormentaTagNoIndex {

			switch fieldType.Type.Kind() {
			case reflect.Slice:
				if err := setMultipleIndexes(txn, v.Field(i), keyRoot, id, fieldType.Name); err != nil {
					return err
				}

			case reflect.Array:
				if err := setMultipleIndexes(txn, v.Field(i), keyRoot, id, fieldType.Name); err != nil {
					return err
				}

			case reflect.String:
				// If the string is tagged with 'split',
				// then index each of the words separately
				if idx == tormentaTagSplit {
					if err := setSplitStringIndexes(txn, v.Field(i), keyRoot, id, fieldType.Name); err != nil {
						return err
					}
				} else {
					if err := setIndexKey(txn, keyRoot, id, fieldType.Name, v.Field(i).Interface()); err != nil {
						return err
					}
				}

			case reflect.Struct:
				// Recursively index embedded structs
				if fieldType.Anonymous {
					indexStruct(v.Field(i), txn, entity, keyRoot, id)
				}

			default:
				if err := setIndexKey(txn, keyRoot, id, fieldType.Name, v.Field(i).Interface()); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func setSplitStringIndexes(txn *badger.Txn, v reflect.Value, root []byte, id gouuidv6.UUID, indexName string) error {
	strings := strings.Split(v.String(), " ")

	// Clean non-content words
	strings = removeNonContentWords(strings)

	for _, s := range strings {
		if err := setIndexKey(txn, root, id, indexName, s); err != nil {
			return err
		}
	}

	return nil
}

func setMultipleIndexes(txn *badger.Txn, v reflect.Value, root []byte, id gouuidv6.UUID, indexName string) error {
	for i := 0; i < v.Len(); i++ {
		if err := setIndexKey(txn, root, id, indexName, v.Index(i).Interface()); err != nil {
			return err
		}
	}

	return nil
}

func setIndexKey(txn *badger.Txn, root []byte, id gouuidv6.UUID, indexName string, indexContent interface{}) error {
	key := makeIndexKey(root, id, indexName, indexContent)

	// Set the index key with no content
	return txn.Set(key, []byte{})
}

// IndexKey returns the bytes of an index key constructed for a particular root, id, index name and index content
func IndexKey(root []byte, id gouuidv6.UUID, indexName string, indexContent interface{}) []byte {
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
