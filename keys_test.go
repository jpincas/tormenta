package tormenta

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/jpincas/gouuidv6"
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
		{"*", ""},
		{"*Test", "test"},
		{"*test.Test", "test"},
		{"*Test.Test.test", "test"},
		{"Test", "test"},
		{"Test.test", "test"},
		{"[]test.test.Test", "test"},
		{"[]*Test.Test.test", "test"},
		{"[]Test", "test"},
		{"[]Test.test", "test"},
		{"[]test.test.Test", "test"},
	}

	for _, test := range testCases {
		result := typeToKeyRoot(test.source)
		if string(result) != test.expectedResult {
			t.Errorf("Converting type sig '%s' to key root produced '%s' instead of '%s'", test.source, result, test.expectedResult)
		}
	}
}

func Test_makeContentKey(t *testing.T) {
	id := newID()

	testCases := []struct {
		testName  string
		root      []byte
		includeID bool
		id        gouuidv6.UUID
		expected  []byte
	}{
		{"No ID", []byte("myentity"), false, id, []byte("c" + keySeparator + "myentity")},
	}

	for _, testCase := range testCases {
		var result []byte

		if testCase.includeID {
			result = newContentKey(testCase.root, testCase.id).bytes()
		} else {
			result = newContentKey(testCase.root).bytes()
		}

		if string(result) != string(testCase.expected) {
			t.Errorf("Testing content key construction (%s). Expecting %s, got %s", testCase.testName, testCase.expected, result)
		}
	}
}

func Test_makeIndexKey(t *testing.T) {

	id := newID()
	ikey := []byte(indexKeyPrefix)

	floatBuf := new(bytes.Buffer)
	var float = 3.14
	binary.Write(floatBuf, binary.LittleEndian, float)

	intBuf := new(bytes.Buffer)
	var int = 3
	binary.Write(intBuf, binary.LittleEndian, uint32(int))

	testCases := []struct {
		testName     string
		root         []byte
		id           gouuidv6.UUID
		indexName    string
		indexContent interface{}
		expected     []byte
	}{
		{
			"no index content",
			[]byte("root"), id, "myindex", nil,
			bytes.Join([][]byte{ikey, []byte("root"), []byte("myindex"), []byte{}, id.Bytes()}, []byte(keySeparator)),
		},
		{
			"string index content",
			[]byte("root"), id, "myindex", "indexContent",
			bytes.Join([][]byte{ikey, []byte("root"), []byte("myindex"), []byte("indexContent"), id.Bytes()}, []byte(keySeparator)),
		},
		{
			"float index content",
			[]byte("root"), id, "myindex", 3.14,
			bytes.Join([][]byte{ikey, []byte("root"), []byte("myindex"), interfaceToBytes(3.14), id.Bytes()}, []byte(keySeparator)),
		},
		{
			"int index content",
			[]byte("root"), id, "myindex", 3,
			bytes.Join([][]byte{ikey, []byte("root"), []byte("myindex"), interfaceToBytes(3), id.Bytes()}, []byte(keySeparator)),
		},
	}

	for _, testCase := range testCases {
		result := makeIndexKey(testCase.root, id, testCase.indexName, testCase.indexContent)
		a := string(result)
		b := string(testCase.expected)
		if a != b {
			t.Errorf("Testing make index key with %s - expected %s, got %s", testCase.testName, b, a)
		}
	}

}
