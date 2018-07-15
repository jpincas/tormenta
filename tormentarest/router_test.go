package tormentadbrest

import (
	"fmt"
	"testing"

	"github.com/jpincas/tormenta"
)

func Test_Router(t *testing.T) {
	order := tormenta.Order{}

	r := makeRouter(&order)
	fmt.Println(r.Routes())
	t.Fail()
}
