package tormenta

import (
	"bytes"
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
