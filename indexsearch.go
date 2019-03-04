package tormenta

import (
	"reflect"

	"github.com/dgraph-io/badger"
	"github.com/jpincas/gouuidv6"
)

type indexSearch struct {
	// Ranges and comparision key
	seekFrom, validTo []byte

	// Reverse fullStruct of searching and returned results
	reverse bool

	// Limit number of returned results
	limit int

	// The entity type being searched
	keyRoot []byte

	// index name
	indexName []byte

	indexKind reflect.Kind

	// Offet - start returning results N entities from the beginning
	// offsetCounter used to track the offset
	offset, offsetCounter int

	// The IDs that we are going to search for in the index
	idsToSearchFor idList

	sumIndexName []byte
	sumTarget    interface{}
}

func (i indexSearch) isLimitMet(noIDsSoFar int) bool {
	return i.limit > 0 && noIDsSoFar >= i.limit
}

func (i *indexSearch) setRanges() {
	i.seekFrom = newIndexKey(i.keyRoot, i.indexName, nil).bytes()
	i.validTo = newIndexKey(i.keyRoot, i.indexName, nil).bytes()

	// For reverse queries, append the byte 0xFF to get inclusive results
	// See Badger issue: https://github.com/dgraph-io/badger/issues/347
	// We can now mark the query as 'reverse prepared'
	if i.reverse {
		i.seekFrom = append(i.seekFrom, 0xFF)
	}
}

func (i indexSearch) getIteratorOptions() badger.IteratorOptions {
	options := badger.DefaultIteratorOptions
	options.Reverse = i.reverse
	options.PrefetchValues = false
	return options
}

func (i indexSearch) execute(txn *badger.Txn) (ids idList) {
	// Set ranges and init the offset counter
	i.setRanges()
	i.offsetCounter = i.offset

	// Create a map of the ids we are looking for
	sourceIDs := map[gouuidv6.UUID]bool{}
	for _, id := range i.idsToSearchFor {
		sourceIDs[id] = true
	}

	it := txn.NewIterator(i.getIteratorOptions())
	defer it.Close()

	for it.Seek(i.seekFrom); it.ValidForPrefix(i.validTo) && !i.isLimitMet(len(ids)); it.Next() {
		item := it.Item()
		thisID := extractID(item.Key())

		// Check to see if this is one of the ids we are looking for.
		// If it is not, continue iterating the index
		if _, ok := sourceIDs[thisID]; !ok {
			continue
		}

		// Skip the first N entities according to the specified offset
		if i.offsetCounter > 0 {
			i.offsetCounter--
			continue
		}

		// If required, take advantage to tally the sum
		if len(i.sumIndexName) > 0 && i.sumTarget != nil {
			quickSum(i.sumTarget, item)
		}

		ids = append(ids, thisID)
	}

	return
}
