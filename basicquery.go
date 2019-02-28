package tormenta

import (
	"time"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

type basicQuery struct {
	// From and To dates
	from, to gouuidv6.UUID

	// Reverse?
	reverse bool

	// Name of the entity -> key root
	keyRoot []byte

	// Limit number of returned results
	limit int

	// Offet - start returning results N entities from the beginning
	// offsetCounter used to track the offset
	offset, offsetCounter int

	// Ranges and comparision key
	seekFrom, validTo, compareTo []byte

	// Is already prepared?
	prepared bool
}

func (b *basicQuery) prepare() {
	b.setFromToIfEmpty()
	b.setRanges()

	// Mark as prepared
	b.prepared = true
}

func (b *basicQuery) setFromToIfEmpty() {
	t1 := time.Time{}
	t2 := time.Now()

	if b.from.IsNil() {
		b.from = fromUUID(t1)
	}

	if b.to.IsNil() {
		b.to = toUUID(t2)
	}
}

func (b *basicQuery) setRanges() {
	var seekFrom, validTo, compareTo []byte

	// For reverse queries, flick-flack start/end and from/to
	// to provide a standardised user API
	if b.reverse {
		tempTo := b.to
		b.to = b.from
		b.from = tempTo
	}

	seekFrom = newContentKey(b.keyRoot, b.from).bytes()
	validTo = newContentKey(b.keyRoot).bytes()
	compareTo = newContentKey(b.keyRoot, b.to).bytes()

	// For reverse queries, append the byte 0xFF to get inclusive results
	// See Badger issue: https://github.com/dgraph-io/badger/issues/347
	if b.reverse {
		seekFrom = append(seekFrom, 0xFF)
	}

	b.seekFrom = seekFrom
	b.validTo = validTo
	b.compareTo = compareTo
}

func (b *basicQuery) reset() {
	b.offsetCounter = b.offset
}

func (b basicQuery) getIteratorOptions() badger.IteratorOptions {
	options := badger.DefaultIteratorOptions
	options.Reverse = b.reverse
	options.PrefetchValues = false
	return options
}

func (b basicQuery) endIteration(it *badger.Iterator, noIDsSoFar int) bool {
	if it.ValidForPrefix(b.validTo) {
		if b.isLimitMet(noIDsSoFar) {
			return false
		}

		return true
	}

	return false
}

func (b basicQuery) isLimitMet(noIDsSoFar int) bool {
	return b.limit > 0 && noIDsSoFar >= b.limit
}

func (b *basicQuery) queryIDs(txn *badger.Txn) (ids idList) {
	if !b.prepared {
		b.prepare()
	}

	b.reset()

	it := txn.NewIterator(b.getIteratorOptions())
	defer it.Close()

	for it.Seek(b.seekFrom); b.endIteration(it, len(ids)); it.Next() {
		// Skip the first N entities according to the specified offset
		if b.offsetCounter > 0 {
			b.offsetCounter--
			continue
		}

		item := it.Item()
		ids = append(ids, extractID(item.Key()))
	}

	return
}
