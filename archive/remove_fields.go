package archive

import (
	"reflect"
	"strings"

	"github.com/buger/jsonparser"
)

// removeSkippedFields will remove from the output JSON any fields that
// have been marked with `tormenta:"-"`
func removeSkippedFields(entityValue reflect.Value, json []byte) {
	for i := 0; i < entityValue.NumField(); i++ {
		fieldType := entityValue.Type().Field(i)
		if shouldDelete, jsonFieldName := shouldDeleteField(fieldType); shouldDelete {
			// TODO: doesnt work with std lib encoded JSON
			jsonparser.Delete(json, jsonFieldName)
		}
	}
}

// shouldDeleteField specifies whether we should delete this field
// from the marshalled JSON output
// according to the optional `tormenta:"_"` tag
func shouldDeleteField(field reflect.StructField) (bool, string) {
	if isTaggedWith(field, tormentaTagNoSave) {
		return getJsonOpts(field)
	}

	return false, ""
}

// Json tags

func getJsonOpts(field reflect.StructField) (bool, string) {
	jsonTag := field.Tag.Get("json")

	// If there is no Json flag, then its a simple delete
	// with the default fieldname
	if jsonTag == "" {
		return true, field.Name
	}

	// Check the options - if the field has been Json tagged
	// with "-" then it won't be in the marshalled Json output
	// anyway, so there's no point trying to delete it
	if jsonTag == "-" {
		return false, ""
	}

	// If there is a Json flag, parse it with the code from
	// the std lib
	overridenFieldName, _ := parseTag(jsonTag)

	// IF we are here then we are good to delete the field
	// we just need to decide whether to use an overriden field name or not
	if overridenFieldName != "" {
		return true, overridenFieldName
	}

	return true, field.Name
}

// This code is copy pasted from the std lib
// so that we deal with JSON tags correctly.
// Here's an explanation of how the std lib deals with JSON tags

// The encoding of each struct field can be customized by the format string
// stored under the "json" key in the struct field's tag.
// The format string gives the name of the field, possibly followed by a
// comma-separated list of options. The name may be empty in order to
// specify options without overriding the default field name.
//
// The "omitempty" option specifies that the field should be omitted
// from the encoding if the field has an empty value, defined as
// false, 0, a nil pointer, a nil interface value, and any empty array,
// slice, map, or string.
//
// As a special case, if the field tag is "-", the field is always omitted.
// Note that a field with name "-" can still be generated using the tag "-,".
//
// Examples of struct field tags and their meanings:
//
//   // Field appears in JSON as key "myName".
//   Field int `json:"myName"`
//
//   // Field appears in JSON as key "myName" and
//   // the field is omitted from the object if its value is empty,
//   // as defined above.
//   Field int `json:"myName,omitempty"`
//
//   // Field appears in JSON as key "Field" (the default), but
//   // the field is skipped if empty.
//   // Note the leading comma.
//   Field int `json:",omitempty"`
//
//   // Field is ignored by this package.
//   Field int `json:"-"`
//
//   // Field appears in JSON as key "-".
//   Field int `json:"-,"`

// https://golang.org/src/encoding/json/tags.go

// tagOptions is the string following a comma in a struct field's "json"
// tag, or the empty string. It does not include the leading comma.
type tagOptions string

// parseTag splits a struct field's json tag into its name and
// comma-separated options.
func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

// Contains reports whether a comma-separated list of options
// contains a particular substr flag. substr must be surrounded by a
// string boundary or commas.
func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}
