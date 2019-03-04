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

func makeIndexKey(root []byte, id gouuidv6.UUID, indexName []byte, indexContent interface{}) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(indexKeyPrefix),
			root,
			indexName,
			interfaceToBytes(indexContent),
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

// interfaceToBytes encodes values to bytes where the underlying interface is the same as the one we want to encode to.  This is used for indexing struct field values where the interface is taken straight from the field value.  The only 'manipulation' required is to cast variable length ints and uints to 32bit length.
func interfaceToBytes(value interface{}) []byte {
	if value == nil {
		return []byte{}
	}

	buf := new(bytes.Buffer)

	switch reflect.ValueOf(value).Type().Kind() {
	case reflect.Int:
		b, _ := interfaceToInt64(value)
		binary.Write(buf, binary.BigEndian, int32(b))
		return flipInt(buf.Bytes())

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		binary.Write(buf, binary.BigEndian, value)
		return flipInt(buf.Bytes())

	case reflect.Float64, reflect.Float32:
		binary.Write(buf, binary.BigEndian, value)
		return flipFloat(buf.Bytes())

	case reflect.Uint:
		b, _ := interfaceToUint64(value)
		binary.Write(buf, binary.BigEndian, uint32(b))
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

// interfaceToBytes encodes values to bytes where the input interface could potentially be anything.  This is used when consulting indices, rather than building them, and thus the input value is provided by the user.  The technique we use is to work out the target type from the original struct field and provide it as an argument to this function.  We then do our best to cast the provided value into the target type.  The best way I've found to do this is for the number types is to output a string and then parse it as a 64bit int, uint or float and then cast to the target numerical type as required.  For bools, if the provided interface is not a bool, we try decoding it as a string (e.g. "false", "f") and then just defualt to False.
func interfaceToBytesWithOverride(value interface{}, typeOverride reflect.Kind) ([]byte, error) {
	if value == nil {
		return []byte{}, nil
	}

	buf := new(bytes.Buffer)

	switch typeOverride {
	// Int
	case reflect.Int:
		b, err := interfaceToInt64(value)
		binary.Write(buf, binary.BigEndian, int32(b))
		return flipInt(buf.Bytes()), err

	case reflect.Int8:
		b, err := interfaceToInt64(value)
		binary.Write(buf, binary.BigEndian, int8(b))
		return flipInt(buf.Bytes()), err

	case reflect.Int16:
		b, err := interfaceToInt64(value)
		binary.Write(buf, binary.BigEndian, int16(b))
		return flipInt(buf.Bytes()), err

	case reflect.Int32:
		b, err := interfaceToInt64(value)
		binary.Write(buf, binary.BigEndian, int32(b))
		return flipInt(buf.Bytes()), err

	case reflect.Int64:
		b, err := interfaceToInt64(value)
		binary.Write(buf, binary.BigEndian, b)
		return flipInt(buf.Bytes()), err

	// Uint
	case reflect.Uint:
		b, err := interfaceToUint64(value)
		binary.Write(buf, binary.BigEndian, uint32(b))
		return buf.Bytes(), err

	case reflect.Uint8:
		b, err := interfaceToUint64(value)
		binary.Write(buf, binary.BigEndian, uint8(b))
		return buf.Bytes(), err

	case reflect.Uint16:
		b, err := interfaceToUint64(value)
		binary.Write(buf, binary.BigEndian, uint16(b))
		return buf.Bytes(), err

	case reflect.Uint32:
		b, err := interfaceToUint64(value)
		binary.Write(buf, binary.BigEndian, uint32(b))
		return buf.Bytes(), err

	case reflect.Uint64:
		b, err := interfaceToUint64(value)
		binary.Write(buf, binary.BigEndian, b)
		return buf.Bytes(), err

	// Float
	case reflect.Float32:
		b, err := interfaceToFloat64(value)
		binary.Write(buf, binary.BigEndian, float32(b))
		return flipFloat(buf.Bytes()), err

	case reflect.Float64:
		b, err := interfaceToFloat64(value)
		binary.Write(buf, binary.BigEndian, b)
		return flipFloat(buf.Bytes()), err

	case reflect.Bool:
		b, err := interfaceToBool(value)
		binary.Write(buf, binary.BigEndian, b)
		return buf.Bytes(), err

	case reflect.Struct:
		// time.Time is a struct, so we encode/decode as int64 (unix seconds)
		if t, ok := reflect.ValueOf(value).Interface().(time.Time); ok {
			binary.Write(buf, binary.BigEndian, t.Unix())
			return flipInt(buf.Bytes()), nil
		}
	}

	// Everything else as a string (lower case)
	return []byte(strings.ToLower(fmt.Sprintf("%v", value))), nil
}

// BIT ORDERING HELPERS

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

// CONVERSIONS

func interfaceToInt64(i interface{}) (int64, error) {
	return strconv.ParseInt(fmt.Sprint(i), 10, 64)
}

func interfaceToUint64(i interface{}) (uint64, error) {
	return strconv.ParseUint(fmt.Sprint(i), 10, 64)
}

func interfaceToFloat64(i interface{}) (float64, error) {
	return strconv.ParseFloat(fmt.Sprint(i), 10)
}

func interfaceToBool(i interface{}) (bool, error) {
	switch i.(type) {
	case bool:
		return i.(bool), nil

	case string:
		s := strings.ToLower(i.(string))
		if s == "true" || s == "t" {
			return true, nil
		} else if s == "false" || s == "f" {
			return false, nil
		}

		return false, fmt.Errorf(ErrIndexTypeBool, s)
	}

	return false, fmt.Errorf(ErrIndexTypeBool, i)
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
