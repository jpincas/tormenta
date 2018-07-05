package tormenta

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/jpincas/gouuidv6"
)

func typeToKeyRoot(typeSig string) []byte {

	sp := strings.Split(typeSig, ".")
	s := sp[len(sp)-1]
	s = strings.TrimPrefix(s, "*")
	s = strings.TrimPrefix(s, "[]")
	s = strings.ToLower(s)

	return []byte(s)
}

const (
	contentKeyPrefix = "c"
	indexKeyPrefix   = "i"
)

func makeKey(root []byte, id gouuidv6.UUID) []byte {
	idBytes := id.Bytes()
	c := []byte(contentKeyPrefix)
	return bytes.Join([][]byte{c, root, idBytes}, []byte(":"))
}

func makePrefix(root, slug []byte) []byte {
	c := []byte(contentKeyPrefix)
	return bytes.Join([][]byte{c, root, slug}, []byte(":"))
}

func makeIndexPrefix(root, indexName []byte, indexValue interface{}) []byte {
	// i:order:customer:5
	i := []byte(indexKeyPrefix)
	return bytes.Join([][]byte{i, root, indexName, interfaceToBytes(indexValue)}, []byte(":"))

}

func getKeyRoot(t interface{}) ([]byte, reflect.Value) {
	e := reflect.Indirect(reflect.ValueOf(t))
	return typeToKeyRoot(e.Type().String()), e
}

func compareKey(id gouuidv6.UUID, root []byte) []byte {
	return bytes.Join([][]byte{root, id.Bytes()}, []byte{})
}

func compareIndexKey(i interface{}, root []byte) []byte {
	return bytes.Join([][]byte{root, interfaceToBytes(i)}, []byte{})
}

func compare(a, b []byte, reverse bool) bool {
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
