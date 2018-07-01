package tormenta

import (
	"bytes"
	"reflect"
	"strings"

	"github.com/jpincas/gouuidv6"
)

func typeToKeyRoot(typeSig string) []byte {
	// Strip off the * (pointer), remove namespace prefixes
	s := strings.TrimPrefix(typeSig, "*")
	sp := strings.Split(s, ".")
	return []byte(sp[len(sp)-1])
}

func makeKey(root []byte, id gouuidv6.UUID) []byte {
	idBytes := id.Bytes()
	return bytes.Join([][]byte{root, idBytes}, []byte(":"))
}

func getKeyRoot(t tormentable) ([]byte, reflect.Value) {
	e := reflect.Indirect(reflect.ValueOf(t))
	return typeToKeyRoot(e.Type().String()), e
}
