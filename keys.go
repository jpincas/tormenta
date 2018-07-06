package tormenta

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/jpincas/gouuidv6"
)

const (
	contentKeyPrefix = "c"
	indexKeyPrefix   = "i"
	keySeparator     = ":"
)

type key struct {
	isIndex      bool
	entityType   []byte
	id           gouuidv6.UUID
	indexName    []byte
	indexContent interface{}
}

// newContentKey returns a key of the correct type
func newContentKey(root []byte, id ...gouuidv6.UUID) key {
	k := key{
		isIndex:    false,
		entityType: root,
	}

	// If an ID is specified
	if len(id) > 0 {
		k.id = id[0]
	}

	return k
}

func newIndexKey(root, indexName []byte, indexContent interface{}) key {
	k := key{
		isIndex:      true,
		entityType:   root,
		indexName:    indexName,
		indexContent: indexContent,
	}

	return k
}

func (k key) bytes() []byte {
	// Use either content/index key prefix
	identifierPrefix := []byte(contentKeyPrefix)
	if k.isIndex {
		identifierPrefix = []byte(indexKeyPrefix)
	}

	// Always start with identifier prefix and entity type
	toJoin := [][]byte{identifierPrefix, k.entityType}

	if k.isIndex {
		// For index keys, now append index name and content
		toJoin = append(toJoin, k.indexName, interfaceToBytes(k.indexContent))
	} else {
		// For content keys, append ID
		// If the ID is nil (i.e. hasn't been added to the struct),
		// then use an empty byteslice, rather than a slice full of zero bytes
		// (which is what id.isNil() means)
		idBytes := k.id.Bytes()
		if k.id.IsNil() {
			idBytes = []byte{}
		}
		toJoin = append(toJoin, idBytes)
	}

	return bytes.Join(toJoin, []byte(keySeparator))
}

// compare compares two key-byte slices
func compareKeyBytes(a, b []byte, reverse bool) bool {
	var r int

	if !reverse {
		r = bytes.Compare(a, b)
	} else {
		r = bytes.Compare(b, a)
	}

	if r > 0 {
		return true
	}

	return false
}

// Key construction helpers

func entityTypeAndValue(t interface{}) ([]byte, reflect.Value) {
	e := reflect.Indirect(reflect.ValueOf(t))
	return typeToKeyRoot(e.Type().String()), e
}

func typeToKeyRoot(typeSig string) []byte {
	sp := strings.Split(typeSig, ".")
	s := sp[len(sp)-1]
	s = strings.TrimPrefix(s, "*")
	s = strings.TrimPrefix(s, "[]")
	s = strings.ToLower(s)

	return []byte(s)
}
