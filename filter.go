package tormenta

import (
	"time"

	"github.com/dgraph-io/badger"

	"github.com/jpincas/gouuidv6"
)

type filter struct {
	////////////////////////////
	// Copied from main query //
	////////////////////////////

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

	/////////////////////////////
	// Specific to this filter //
	/////////////////////////////

	// The start and end points of the index range search
	start, end interface{}

	// Name of the index on which to apply filter
	indexName []byte

	// Is this a 'starts with' index query
	isStartsWithQuery bool

	// Ranges and comparision key
	seekFrom, validTo, compareTo []byte

	// Is already prepared?
	prepared bool
}

func (f filter) isIndexRangeSearch() bool {
	return f.start != f.end && !f.isStartsWithQuery
}

func (f filter) isExactIndexMatchSearch() bool {
	return f.start == f.end && f.start != nil && f.end != nil
}

func (f *filter) prepare() {
	// 'starts with' type query doesn't work with reverse
	// so switch it back to a regular search
	if f.isStartsWithQuery && f.reverse {
		f.reverse = false
	}

	f.setFromToIfEmpty()
	f.setRanges()

	// Mark as prepared
	f.prepared = true
}

func (f *filter) setFromToIfEmpty() {

	// For index range searches - we don't do this, so exit right away
	if f.isIndexRangeSearch() {
		return
	}

	// If 'from' or 'to' have not been specified manually by the user,
	// then we set them to the 'widest' times possible,
	// i.e. 'between beginning of time' and 'now'
	// If we don't do this, then some searches work OK, but particuarly reversed searches
	// can experience strange behaviour (namely returned 0 results), because the iteration
	// ends up starting from the end of the list.
	// Another side-effect of not doing this is that exact match string searches would become 'starts with' searches.  We might want that behaviour though, so we include a check for this type of search below
	t1 := time.Time{}
	t2 := time.Now()

	if f.from.IsNil() {
		// If we are doing a 'starts with' query,
		// then we DON'T want to set the from point
		// This magically gives us 'starts with'
		// instead of exact match,
		// BUT - this trick only works for forward searches,
		// not 'reverse' searches,
		// so there is a protection in the query preparation
		if !f.isStartsWithQuery {
			f.from = fromUUID(t1)
		}
	}

	if f.to.IsNil() {
		f.to = toUUID(t2)
	}
}

func (f *filter) setRanges() {
	var seekFrom, validTo, compareTo []byte

	// For reverse queries, flick-flack start/end and from/to
	// to provide a standardised user API
	if f.reverse {
		tempEnd := f.end
		f.end = f.start
		f.start = tempEnd

		tempTo := f.to
		f.to = f.from
		f.from = tempTo
	}

	if f.isExactIndexMatchSearch() {
		// For index searches with exact match
		seekFrom = newIndexMatchKey(f.keyRoot, f.indexName, f.start, f.from).bytes()
		validTo = newIndexMatchKey(f.keyRoot, f.indexName, f.end).bytes()
		compareTo = newIndexMatchKey(f.keyRoot, f.indexName, f.end, f.to).bytes()
	} else {
		// For regular index searches
		seekFrom = newIndexKey(f.keyRoot, f.indexName, f.start).bytes()
		validTo = newIndexKey(f.keyRoot, f.indexName, nil).bytes()
		compareTo = newIndexKey(f.keyRoot, f.indexName, f.end).bytes()
	}

	// For reverse queries, append the byte 0xFF to get inclusive results
	// See Badger issue: https://github.com/dgraph-io/badger/issues/347
	if f.reverse {
		seekFrom = append(seekFrom, 0xFF)
	}

	f.seekFrom = seekFrom
	f.validTo = validTo
	f.compareTo = compareTo
}

func (f filter) endIteration(it *badger.Iterator, noIDsSoFar int) bool {
	if it.ValidForPrefix(f.validTo) {
		if f.isLimitMet(noIDsSoFar) || f.isEndOfRange(it) {
			return false
		}

		return true
	}

	return false
}

func (f filter) shouldStripKeyID() bool {
	// Index queries which are exact match AND have a 'to' clause
	// also never need to have ID stripped
	if f.isExactIndexMatchSearch() && !f.to.IsNil() {
		return false
	}

	return true
}

func (f filter) isEndOfRange(it *badger.Iterator) bool {
	key := it.Item().Key()
	return f.end != nil && compareKeyBytes(f.compareTo, key, f.reverse, f.shouldStripKeyID())
}

func (f filter) isLimitMet(noIDsSoFar int) bool {
	return f.limit > 0 && noIDsSoFar >= f.limit
}

func (f *filter) reset() {
	f.offsetCounter = f.offset
}

func (f filter) getIteratorOptions() badger.IteratorOptions {
	options := badger.DefaultIteratorOptions
	options.Reverse = f.reverse
	options.PrefetchValues = false
	return options
}

func (f *filter) queryIDs(txn *badger.Txn) (ids idList) {
	if !f.prepared {
		f.prepare()
	}

	f.reset()

	it := txn.NewIterator(f.getIteratorOptions())
	defer it.Close()

	for it.Seek(f.seekFrom); f.endIteration(it, ids.length()); it.Next() {
		// If this is a 'range index' type Query
		// that ALSO has a date range, the procedure is a little more complicated
		// compared to an exact index match.
		// Since the start/end points of the iteration focus on the index, e.g. E-J (alphabetical index)
		// we need to manually check all the keys and reject those that don't fit the date range
		if !f.isExactIndexMatchSearch() {
			key := extractID(it.Item().Key())
			if keyIsOutsideDateRange(key, f.from, f.to) {
				continue
			}
		}

		// Skip the first N entities according to the specified offset
		if f.offsetCounter > 0 {
			f.offsetCounter--
			continue
		}

		item := it.Item()
		ids = append(ids, extractID(item.Key()))
	}

	return
}

// Helpers

func toIndexName(s string) []byte {
	return []byte(s)
}
