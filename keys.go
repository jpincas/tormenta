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

func makeKey(root []byte, id gouuidv6.UUID) []byte {
	idBytes := id.Bytes()
	return bytes.Join([][]byte{root, idBytes}, []byte(":"))
}

func getKeyRoot(t interface{}) ([]byte, reflect.Value) {
	e := reflect.Indirect(reflect.ValueOf(t))
	return typeToKeyRoot(e.Type().String()), e
}
