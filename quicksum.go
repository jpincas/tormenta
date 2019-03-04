package tormenta

import "github.com/dgraph-io/badger"

func quickSum(target interface{}, item *badger.Item) {
	// TODO: is there a more efficient way to increment
	// the sum target given that we don't know what type it is
	switch target.(type) {

	// Signed Ints
	case *int:
		// Reminder - decoding the index values only works for fixed length integers
		// So in the case of ints, we set up an int32 target and use
		// that to accumulate
		acc := *target.(*int)
		var int32target int32
		extractIndexValue(item.Key(), &int32target)
		*target.(*int) = acc + int(int32target)
	case *int8:
		acc := *target.(*int8)
		extractIndexValue(item.Key(), target)
		*target.(*int8) = acc + *target.(*int8)
	case *int16:
		acc := *target.(*int16)
		extractIndexValue(item.Key(), target)
		*target.(*int16) = acc + *target.(*int16)
	case *int32:
		acc := *target.(*int32)
		extractIndexValue(item.Key(), target)
		*target.(*int32) = acc + *target.(*int32)
	case *int64:
		acc := *target.(*int64)
		extractIndexValue(item.Key(), target)
		*target.(*int64) = acc + *target.(*int64)

	// Unsigned ints
	case *uint:
		// See above for notes on variable vs fixed length
		acc := *target.(*uint)
		var uint32target uint32
		extractIndexValue(item.Key(), &uint32target)
		*target.(*uint) = acc + uint(uint32target)
	case *uint8:
		acc := *target.(*uint8)
		extractIndexValue(item.Key(), target)
		*target.(*uint8) = acc + *target.(*uint8)
	case *uint16:
		acc := *target.(*uint16)
		extractIndexValue(item.Key(), target)
		*target.(*uint16) = acc + *target.(*uint16)
	case *uint32:
		acc := *target.(*uint32)
		extractIndexValue(item.Key(), target)
		*target.(*uint32) = acc + *target.(*uint32)
	case *uint64:
		acc := *target.(*uint64)
		extractIndexValue(item.Key(), target)
		*target.(*uint64) = acc + *target.(*uint64)

	// Floats
	case *float64:
		acc := *target.(*float64)
		extractIndexValue(item.Key(), target)
		*target.(*float64) = acc + *target.(*float64)

	case *float32:
		acc := *target.(*float32)
		extractIndexValue(item.Key(), target)
		*target.(*float32) = acc + *target.(*float32)
	}
}
