package tormenta

import (
	"reflect"
	"strings"
)

const (
	tormentaTag        = "tormenta"
	tormentaTagNoIndex = "noindex"
	tormentaTagSplit   = "split"
	tagSeparator       = ";"
)

// Tormenta-specific tags

func getTormentaTags(field reflect.StructField) []string {
	compositeTag := field.Tag.Get(tormentaTag)
	return strings.Split(compositeTag, tagSeparator)
}

func isTaggedWith(field reflect.StructField, targetTag string) bool {
	tags := getTormentaTags(field)
	for _, tag := range tags {
		if tag == targetTag {
			return true
		}
	}

	return false
}

// shouldIndex specifies whether a field should be indexed or not
// according to the optional `tormenta:"noindex"` tag
func shouldIndex(field reflect.StructField) bool {
	return !isTaggedWith(field, tormentaTagNoIndex)
}
