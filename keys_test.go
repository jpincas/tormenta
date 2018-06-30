package tormenta

import (
	"testing"
)

func Test_typeToKeyRoot(t *testing.T) {
	testCases := []struct {
		source         string
		expectedResult string
	}{
		{"*", ""},
		{"*test", "test"},
		{"*test.test", "test"},
		{"*test.test.test", "test"},
		{"test", "test"},
		{"test.test", "test"},
		{"test.test.test", "test"},
	}

	for _, test := range testCases {
		result := typeToKeyRoot(test.source)
		if string(result) != test.expectedResult {
			t.Errorf("Converting type sig '%s' to key root produced '%s' instead of '%s'", test.source, result, test.expectedResult)
		}
	}
}
