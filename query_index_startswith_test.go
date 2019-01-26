package tormenta_test

import (
	"errors"
	"testing"

	"github.com/jpincas/tormenta"
	"github.com/jpincas/tormenta/demo"
)

func Test_IndexQuery_StartsWith(t *testing.T) {
	customers := []string{"j", "jo", "jon", "jonathan", "job", "pablo"}
	var orders []tormenta.Tormentable

	for _, customer := range customers {
		orders = append(orders, &demo.Order{
			Customer: customer,
		})
	}

	db, _ := tormenta.OpenTest("data/tests")
	defer db.Close()
	db.Save(orders...)

	testCases := []struct {
		testName      string
		startsWith    string
		reverse       bool
		expected      int
		expectedError error
	}{
		{"blank string", "", false, 0, errors.New(tormenta.ErrBlankInputStartsWithQuery)},
		{"no match - no interference", "nocustomerwiththisname", false, 0, nil},
		{"single match - no interference", "pablo", false, 1, nil},
		{"single match - possible interference", "jonathan", false, 1, nil},
		{"single match - possible interference", "job", false, 1, nil},
		{"wide match - 1 letter", "j", false, 5, nil},
		{"wide match - 2 letters", "jo", false, 4, nil},
		{"wide match - 3 letters", "jon", false, 2, nil},

		// Reversed - shouldn't make any difference to N
		{"blank string", "", true, 0, errors.New(tormenta.ErrBlankInputStartsWithQuery)},
		{"no match - no interference", "nocustomerwiththisname", true, 0, nil},
		{"single match - no interference", "pablo", true, 1, nil},
		{"single match - possible interference", "jonathan", true, 1, nil},
		{"single match - possible interference", "job", true, 1, nil},
		{"wide match - 1 letter", "j", true, 5, nil},
		{"wide match - 2 letters", "jo", true, 4, nil},
		{"wide match - 3 letters", "jon", true, 2, nil},
	}

	for _, testCase := range testCases {
		results := []demo.Order{}

		q := db.Find(&results).StartsWith("customer", testCase.startsWith)
		if testCase.reverse {
			q.Reverse()
		}

		n, _, err := q.Run()

		if testCase.expectedError != nil && err == nil {
			t.Errorf("Testing %s. Expected error [%v] but got none", testCase.testName, testCase.expectedError)
		}

		if testCase.expectedError == nil && err != nil {
			t.Errorf("Testing %s. Didn't expect error [%v]", testCase.testName, err)
		}

		if n != testCase.expected {
			t.Errorf("Testing %s.  Expecting %v, got %v", testCase.testName, testCase.expected, n)
		}
	}
}
