package tormenta

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/jpincas/gouuidv6"
)

const (
	contentKeyPrefix = "c"
	indexKeyPrefix   = "i"
	keySeparator     = "~±^"
)

type key struct {
	isIndex      bool
	entityType   []byte
	id           gouuidv6.UUID
	indexName    []byte
	indexContent interface{}
	exactMatch   bool
}

// newContentKey returns a key of the correct type
func newContentKey(root []byte, id ...gouuidv6.UUID) key {
	return withID(key{
		isIndex:    false,
		entityType: root,
	}, id)
}

func newIndexKey(root, indexName []byte, indexContent interface{}, id ...gouuidv6.UUID) key {
	return withID(key{
		isIndex:      true,
		entityType:   root,
		indexName:    indexName,
		indexContent: indexContent,
	}, id)
}

func newIndexMatchKey(root, indexName []byte, indexContent interface{}, id ...gouuidv6.UUID) key {
	return withID(key{
		isIndex:      true,
		exactMatch:   true,
		entityType:   root,
		indexName:    indexName,
		indexContent: indexContent,
	}, id)
}

func withID(k key, id []gouuidv6.UUID) key {
	// If an ID is specified
	if len(id) > 0 {
		k.id = id[0]
	}

	return k
}

func (k key) shouldAppendID() bool {
	// If ID is nil, definite no
	if k.id.IsNil() {
		return false
	}

	// For non-index keys, do append
	if !k.isIndex {
		return true
	}

	// For index keys using exact matching, do append
	if k.exactMatch {
		return true
	}

	return false
}

// c:fullStructs:sdfdsf-9sdfsdf-8dsf-sdf-9sdfsdf
// i:fullStructs:department:3
// i:fullStructs:department:3:sdfdsf-9sdfsdf-8dsf-sdf-9sdfsdf

func (k key) bytes() []byte {
	// Use either content/index key prefix
	identifierPrefix := []byte(contentKeyPrefix)
	if k.isIndex {
		identifierPrefix = []byte(indexKeyPrefix)
	}

	// Always start with identifier prefix and entity type
	toJoin := [][]byte{identifierPrefix, k.entityType}

	// For index keys, now append index name and content
	if k.isIndex {
		toJoin = append(toJoin, k.indexName, interfaceToBytes(k.indexContent))
	}

	if k.shouldAppendID() {
		toJoin = append(toJoin, k.id.Bytes())
	}

	return bytes.Join(toJoin, []byte(keySeparator))
}

func extractID(b []byte) (uuid gouuidv6.UUID) {
	s := bytes.Split(b, []byte(keySeparator))
	idBytes := s[len(s)-1]
	copy(uuid[:], idBytes)
	return
}

func extractIndexValue(b []byte, i interface{}) {
	s := bytes.Split(b, []byte(keySeparator))
	indexValueBytes := s[3]

	// For unsigned ints, we need to flip the sign bit back
	switch i.(type) {
	case *int, *int8, *int16, *int32, *int64:
		flipInt(indexValueBytes)
	case *float64, *float32:
		revertFloat(indexValueBytes)
	}

	buf := bytes.NewBuffer(indexValueBytes)
	binary.Read(buf, binary.BigEndian, i) //TODO: error handling
}

func stripID(b []byte) []byte {
	s := bytes.Split(b, []byte(keySeparator))
	return bytes.Join(s[:len(s)-1], []byte(keySeparator))
}

// compare compares two key-byte slices
func compareKeyBytes(a, b []byte, reverse bool, removeID bool) bool {
	if removeID {
		b = stripID(b)
	}

	var r int

	if !reverse {
		r = bytes.Compare(a, b)
	} else {
		r = bytes.Compare(b, a)
	}

	if r < 0 {
		return true
	}

	return false
}

func keyIsOutsideDateRange(key, start, end gouuidv6.UUID) bool {
	// No dates at all? Then its definitely not outside the range
	if start.IsNil() && end.IsNil() {
		return false
	}

	// For start date only
	if end.IsNil() {
		return key.Compare(start)
	}

	// For both start and end
	return key.Compare(start) || !key.Compare(end)
}

// Key construction helpers

func KeyRoot(t interface{}) []byte {
	k, _ := entityTypeAndValue(t)
	return k
}

func KeyRootString(entity Record) string {
	return string(KeyRoot(entity))
}

func typeToKeyRoot(typeSig string) []byte {
	return []byte(strings.ToLower(cleanType(typeSig)))
}

func cleanType(typeSig string) string {
	sp := strings.Split(typeSig, ".")
	s := sp[len(sp)-1]
	s = strings.TrimPrefix(s, "*")
	s = strings.TrimPrefix(s, "[]")

	return s
}
