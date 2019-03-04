package tormenta

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
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
		nil,
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
		nil,
	)

	for i := range keys {
		if err := txn.Delete(keys[i]); err != nil {
			return err
		}
	}

	return nil
}

func indexStruct(v reflect.Value, entity Record, keyRoot []byte, id gouuidv6.UUID, path []byte) (keys [][]byte) {
	for i := 0; i < v.NumField(); i++ {

		fieldType := v.Type().Field(i)
		indexName := []byte(fieldType.Name)
		if path != nil {
			indexName = nestedIndexKeyRoot(path, indexName)
		}

		if !isTaggedWith(fieldType, tormentaTagNoIndex, tormentaTagNoSave) {

			switch fieldType.Type.Kind() {

			// Slice: index members individually
			case reflect.Slice:
				keys = append(keys, getMultipleIndexKeys(v.Field(i), keyRoot, id, indexName)...)

			// Array: index members individually
			case reflect.Array:
				// UUIDV6s are arrays, so we intercept them here
				if fieldType.Type == reflect.TypeOf(gouuidv6.UUID{}) {
					keys = append(keys, makeIndexKey(keyRoot, id, indexName, v.Field(i).Interface()))
				} else {
					keys = append(keys, getMultipleIndexKeys(v.Field(i), keyRoot, id, indexName)...)
				}

			// Strings: either straight index, or split by words
			case reflect.String:
				if isTaggedWith(fieldType, tormentaTagSplit) {
					keys = append(keys, getSplitStringIndexes(v.Field(i), keyRoot, id, indexName)...)
				} else {
					keys = append(keys, makeIndexKey(keyRoot, id, indexName, v.Field(i).Interface()))
				}

			// Anonymous/ Nested Structs
			case reflect.Struct:
				// time.Time is a struct, so we'll intercept it here
				// and send it to the index key maker which will translate it to int64
				// see below interfaceToBytes for more on that
				f := v.Field(i).Interface()
				if _, ok := f.(time.Time); ok {
					keys = append(keys, makeIndexKey(keyRoot, id, indexName, f))
				}

				// Recursively index embedded structs
				if fieldType.Anonymous {
					keys = append(keys, indexStruct(v.Field(i), entity, keyRoot, id, nil)...)
				}

				// And named structs, if they are tagged 'nested'
				// But construct the index with path separators
				if isTaggedWith(fieldType, tormentaTagNestedIndex) {
					keys = append(keys, indexStruct(v.Field(i), entity, keyRoot, id, indexName)...)
				}

			default:
				keys = append(keys, makeIndexKey(keyRoot, id, indexName, v.Field(i).Interface()))
			}
		}
	}

	return
}

// MakeIndexKey constructs an index key
func MakeIndexKey(root []byte, id gouuidv6.UUID, indexName []byte, indexContent interface{}) []byte {
	return makeIndexKey(root, id, indexName, indexContent)
}

func makeIndexKey(root []byte, id gouuidv6.UUID, indexName []byte, indexContent interface{}, typeOverride ...reflect.Kind) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(indexKeyPrefix),
			root,
			indexName,
			interfaceToBytes(indexContent, typeOverride...),
			id.Bytes(),
		},
		[]byte(keySeparator),
	)
}

func getMultipleIndexKeys(v reflect.Value, root []byte, id gouuidv6.UUID, indexName []byte) (keys [][]byte) {
	for i := 0; i < v.Len(); i++ {
		key := makeIndexKey(root, id, indexName, v.Index(i).Interface())
		keys = append(keys, key)
	}

	return
}

func getSplitStringIndexes(v reflect.Value, root []byte, id gouuidv6.UUID, indexName []byte) (keys [][]byte) {
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
func interfaceToBytes(value interface{}, typeOverride ...reflect.Kind) []byte {
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
	if len(typeOverride) == 0 {
		// If a type override has not been specified, then this request is coming from
		// the index construction precedure, where the value for the index is being taken
		// straight of the struct field, therefore it will always be the 'right' type
		switch reflect.ValueOf(value).Type().Kind() {
		case reflect.Int:
			binary.Write(buf, binary.BigEndian, int32(interfaceToInt64(value)))
			return flipInt(buf.Bytes())

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			binary.Write(buf, binary.BigEndian, value)
			return flipInt(buf.Bytes())

		case reflect.Float64, reflect.Float32:
			binary.Write(buf, binary.BigEndian, value)
			return flipFloat(buf.Bytes())

		case reflect.Uint:
			binary.Write(buf, binary.BigEndian, int32(interfaceToUint64(value)))
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
	} else {
		// If a type override has been specified then this request is coming from
		// a query builder, where the supplied index value could potentially be anything.
		// We therefore first have to attempt to cajole the supplied index value into the
		// original type of that field on the struct
		switch typeOverride[0] {
		// Int
		case reflect.Int:
			binary.Write(buf, binary.BigEndian, int32(interfaceToInt64(value)))
			return flipInt(buf.Bytes())

		case reflect.Int8:
			binary.Write(buf, binary.BigEndian, int8(interfaceToInt64(value)))
			return flipInt(buf.Bytes())

		case reflect.Int16:
			binary.Write(buf, binary.BigEndian, int16(interfaceToInt64(value)))
			return flipInt(buf.Bytes())

		case reflect.Int32:
			binary.Write(buf, binary.BigEndian, int32(interfaceToInt64(value)))
			return flipInt(buf.Bytes())

		case reflect.Int64:
			binary.Write(buf, binary.BigEndian, interfaceToInt64(value))
			return flipInt(buf.Bytes())

		// Uint
		case reflect.Uint:
			binary.Write(buf, binary.BigEndian, uint32(interfaceToUint64(value)))
			return buf.Bytes()

		case reflect.Uint8:
			binary.Write(buf, binary.BigEndian, uint8(interfaceToUint64(value)))
			return buf.Bytes()

		case reflect.Uint16:
			binary.Write(buf, binary.BigEndian, uint16(interfaceToUint64(value)))
			return buf.Bytes()

		case reflect.Uint32:
			binary.Write(buf, binary.BigEndian, uint32(interfaceToUint64(value)))
			return buf.Bytes()

		case reflect.Uint64:
			binary.Write(buf, binary.BigEndian, interfaceToUint64(value))
			return buf.Bytes()

		// Float
		case reflect.Float32:
			binary.Write(buf, binary.BigEndian, float32(interfaceToFloat64(value)))
			return flipFloat(buf.Bytes())

		case reflect.Float64:
			binary.Write(buf, binary.BigEndian, interfaceToFloat64(value))
			return flipFloat(buf.Bytes())

		case reflect.Bool:
			binary.Write(buf, binary.BigEndian, interfaceToBool(value))
			return buf.Bytes()

		case reflect.Struct:
			// time.Time is a struct, so we encode/decode as int64 (unix seconds)
			if t, ok := reflect.ValueOf(value).Interface().(time.Time); ok {
				binary.Write(buf, binary.BigEndian, t.Unix())
				return flipInt(buf.Bytes())
			}
		}
	}

	// Everything else as a string (lower case)
	return []byte(strings.ToLower(fmt.Sprintf("%v", value)))
}

func flipInt(b []byte) []byte {
	// Deal with 0
	if len(b) == 0 {
		return b
	}

	b[0] ^= 1 << 7
	return b
}

func flipFloat(b []byte) []byte {
	// Deal with 0
	if len(b) == 0 {
		return b
	}

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
	// Deal with 0
	if len(b) == 0 {
		return b
	}

	if b[0]>>7 > 0 {
		b[0] ^= 0x80
	} else {
		for i, bb := range b {
			b[i] = ^bb
		}
	}

	return b
}

// Conversion Functions

func interfaceToInt64(i interface{}) int64 {
	r, _ := strconv.ParseInt(fmt.Sprint(i), 10, 64)
	return r
}

func interfaceToUint64(i interface{}) uint64 {
	r, _ := strconv.ParseUint(fmt.Sprint(i), 10, 64)
	return r
}

func interfaceToFloat64(i interface{}) float64 {
	r, _ := strconv.ParseFloat(fmt.Sprint(i), 10)
	return r
}

func interfaceToBool(i interface{}) bool {
	switch i.(type) {
	case bool:
		return i.(bool)

	case string:
		s := strings.ToLower(i.(string))
		if s == "true" || s == "t" {
			return true
		}

		return false

	case int:
		if i.(int) > 0 {
			return true
		}

		return false

	case float64:
		if i.(float64) > 0 {
			return true
		}

		return false
	}

	return false
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
