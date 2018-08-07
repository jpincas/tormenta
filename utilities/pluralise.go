package utilities

import (
	"github.com/jinzhu/inflection"
)

func Pluralise(s string) string {
	return inflection.Plural(s)
}
