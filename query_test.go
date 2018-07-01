package tormenta

import (
	"fmt"
	"testing"
)

func Test_BasicQuery(t *testing.T) {
	db, _ := OpenTest("data/tests")
	defer db.Close()

	var orders []Order

	db.Query(&orders).Run()
	fmt.Println(len(orders))
	t.Fail()
}
