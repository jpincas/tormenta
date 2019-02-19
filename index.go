package tormenta

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

// Index format
// i:indexname:root:indexcontent:entityID
// i:fullStruct:customer:5:324ds-3werwf-234wef-23wef

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

func deIndex(txn *badger.Txn, entity Record) error {
	keys := indexStruct(
		recordValue(entity),
		entity,
		KeyRoot(entity),
		entity.GetID(),
	)

	for i := range keys {
		if err := txn.Delete(keys[i]); err != nil {
			return err
		}
	}

	return nil
}

func indexStruct(v reflect.Value, entity Record, keyRoot []byte, id gouuidv6.UUID) (keys [][]byte) {
	for i := 0; i < v.NumField(); i++ {

		// If the 'tormenta:noindex' tag is present, don't index
		fieldType := v.Type().Field(i)
		if !isTaggedWith(fieldType, tormentaTagNoIndex, tormentaTagNoSave) {
			switch fieldType.Type.Kind() {
			case reflect.Slice:
				keys = append(keys, getMultipleIndexKeys(v.Field(i), keyRoot, id, fieldType.Name)...)

			case reflect.Array:
				// UUIDV6s are arrays, so we intercept them here
				if fieldType.Type == reflect.TypeOf(gouuidv6.UUID{}) {
					keys = append(keys, makeIndexKey(keyRoot, id, fieldType.Name, v.Field(i).Interface()))
				} else {
					keys = append(keys, getMultipleIndexKeys(v.Field(i), keyRoot, id, fieldType.Name)...)
				}

			case reflect.String:
				// If the string is tagged with 'split',
				// then index each of the words separately
				if isTaggedWith(fieldType, tormentaTagSplit) {
					keys = append(keys, getSplitStringIndexes(v.Field(i), keyRoot, id, fieldType.Name)...)
				} else {
					keys = append(keys, makeIndexKey(keyRoot, id, fieldType.Name, v.Field(i).Interface()))
				}

			case reflect.Struct:
				// time.Time is a struct, so we'll intercept it here
				// and send it to the index key maker which will translate it to int64
				// see below interfaceToBytes for more on that
				f := v.Field(i).Interface()
				if _, ok := f.(time.Time); ok {
					keys = append(keys, makeIndexKey(keyRoot, id, fieldType.Name, f))
				}

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

// interfaceToBytes encodes values to bytes
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
	case reflect.Int:
		binary.Write(buf, binary.BigEndian, intInterfaceToInt32(value))
		return flipInt(buf.Bytes())

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		binary.Write(buf, binary.BigEndian, value)
		return flipInt(buf.Bytes())

	case reflect.Float64, reflect.Float32:
		binary.Write(buf, binary.BigEndian, value)
		return flipFloat(buf.Bytes())

	case reflect.Uint:
		binary.Write(buf, binary.BigEndian, intInterfaceToUInt32(value))
		return buf.Bytes()

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		binary.Write(buf, binary.BigEndian, value)
		return buf.Bytes()

	case reflect.Bool:
		binary.Write(buf, binary.BigEndian, value)
		return buf.Bytes()

	case reflect.Struct:
		// time.Time is a struct, so we encode/decode as int64 (unix seconds)
		if t, ok := reflect.ValueOf(value).Interface().(time.Time); ok {
			binary.Write(buf, binary.BigEndian, t.Unix())
			return flipInt(buf.Bytes())
		}

	}

	// Everything else as a string (lower case)
	return []byte(strings.ToLower(fmt.Sprintf("%v", value)))
}

func flipInt(b []byte) []byte {
	b[0] ^= 1 << 7
	return b
}

func flipFloat(b []byte) []byte {
	if b[0]>>7 > 0 {
		for i, bb := range b {
			b[i] = ^bb
		}
	} else {
		b[0] ^= 0x80
	}

	return b
}

func revertFloat(b []byte) []byte {
	if b[0]>>7 > 0 {
		b[0] ^= 0x80
	} else {
		for i, bb := range b {
			b[i] = ^bb
		}
	}

	return b
}

// isNegative uses the bitwise &operator to determine if the first bit of a bit slice is 1.
// In the case of signed numbers, this means its a negative
// func isNegative(b []byte) bool {
// 	return b[0]>>7 > 0
// }

// I think a simple modification can extend this to negative numbers: XOR all positive numbers with 0x8000... and negative numbers with 0xffff.... This should flip the sign bit on both (so negative numbers go first), and then reverse the ordering on negative numbers. Does anyone see a problem with this?
// func signByteFloat(b []byte) []byte {
// 	return b
// }
