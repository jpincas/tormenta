package tormenta

import (
	"fmt"
	"testing"
	"time"

	"github.com/jpincas/gouuidv6"
)

var (
	id1 = idFromInt(1)
	id2 = idFromInt(2)
	id3 = idFromInt(3)
	id4 = idFromInt(4)
	id5 = idFromInt(5)
	id6 = idFromInt(6)
	id7 = idFromInt(7)
	id8 = idFromInt(8)
	id9 = idFromInt(9)
)

func Test_Sort(t *testing.T) {

	testCases := []struct {
		testName string
		unsorted idList
		expected idList
		reverse  bool
	}{
		{
			"empty list",
			idList{},
			idList{},
			true,
		},
		{
			"single member",
			idList{id1},
			idList{id1},
			true,
		},
		{
			"multiple members - preserve order",
			idList{id5, id4, id3, id2, id1},
			idList{id5, id4, id3, id2, id1},
			true,
		},
		{
			"multiple members - change order",
			idList{id1, id2, id3, id4, id5},
			idList{id5, id4, id3, id2, id1},
			true,
		},
		{
			"multiple members - change order - oldest first",
			idList{id5, id4, id3, id2, id1},
			idList{id1, id2, id3, id4, id5},
			false,
		},
		{
			"multiple members - preserve order - oldest first",
			idList{id1, id2, id3, id4, id5},
			idList{id1, id2, id3, id4, id5},
			false,
		},
	}

	for _, testCase := range testCases {
		testCase.unsorted.sort(testCase.reverse)
		if err := compareIDLists(testCase.unsorted, testCase.expected); err != nil {
			t.Errorf("Testing: %s. Got error: %v", testCase.testName, err)
		}
	}

}

func Test_Union(t *testing.T) {

	testCases := []struct {
		testName string
		idLists  []idList
		expected idList
	}{
		{
			"empty list",
			[]idList{},
			idList{},
		},
		{
			"1 list (empty)",
			[]idList{idList{}},
			idList{},
		},
		{
			"2 lists (both empty)",
			[]idList{idList{}, idList{}},
			idList{},
		},
		{
			"1 list (1 member)",
			[]idList{
				idList{id1},
			},
			idList{id1},
		},
		{
			"1 list (multiple members, sort not required)",
			[]idList{
				idList{id3, id2, id1},
			},
			idList{id3, id2, id1},
		},
		{
			"1 list (multiple members, sort required)",
			[]idList{
				idList{id1, id2, id3},
			},
			idList{id3, id2, id1},
		},
		{
			"2 lists (multiple members, no overlap)",
			[]idList{
				idList{id3, id2, id1},
				idList{id6, id5, id4},
			},
			idList{id6, id5, id4, id3, id2, id1},
		},
		{
			"2 lists (multiple members,  overlap)",
			[]idList{
				idList{id3, id2, id1},
				idList{id5, id4, id3},
			},
			idList{id5, id4, id3, id2, id1},
		},
		{
			"3 lists (multiple members, overlap, repeats)",
			[]idList{
				idList{id3, id2, id1},
				idList{id5, id4, id3},
				idList{id5, id5, id1},
			},
			idList{id5, id4, id3, id2, id1},
		},
	}

	for _, testCase := range testCases {
		result := union(testCase.idLists...)
		// The expected results implicate a reverse sort of the results -
		// that's just how I wrote them originally
		result.sort(true)
		if err := compareIDLists(result, testCase.expected); err != nil {
			t.Errorf("Testing: %s. Got error: %v", testCase.testName, err)
		}
	}

}

func Test_Intersection(t *testing.T) {

	testCases := []struct {
		testName string
		idLists  []idList
		expected idList
	}{
		{
			"empty list",
			[]idList{},
			idList{},
		},
		{
			"1 list (empty)",
			[]idList{idList{}},
			idList{},
		},
		{
			"2 lists (both empty)",
			[]idList{idList{}, idList{}},
			idList{},
		},
		{
			"1 list (1 member)",
			[]idList{
				idList{id1},
			},
			idList{id1},
		},
		{
			"1 list (multiple members, sort not required)",
			[]idList{
				idList{id3, id2, id1},
			},
			idList{id3, id2, id1},
		},
		{
			"1 list (multiple members, sort required)",
			[]idList{
				idList{id1, id2, id3},
			},
			idList{id3, id2, id1},
		},
		{
			"2 lists (multiple members, no overlap)",
			[]idList{
				idList{id3, id2, id1},
				idList{id6, id5, id4},
			},
			idList{},
		},
		{
			"2 lists (multiple members,  overlap)",
			[]idList{
				idList{id3, id2, id1},
				idList{id5, id4, id3},
			},
			idList{id3},
		},
		{
			"3 lists (multiple members, overlap, repeats)",
			[]idList{
				idList{id3, id5, id3},
				idList{id5, id4, id3},
				idList{id5, id5, id3},
			},
			idList{id5, id3},
		},
		{
			"complete example",
			[]idList{
				idList{id1, id2, id3, id4, id5, id6, id7, id8, id9},
				idList{id5, id4, id3},
				idList{id3, id4, id5, id3, id4},
				idList{id3, id4, id5, id2, id1},
				idList{id5, id4, id3, id7, id8},
			},
			idList{id5, id4, id3},
		},
	}

	for _, testCase := range testCases {
		result := intersection(testCase.idLists...)
		// The expected results implicate a reverse sort of the results -
		// that's just how I wrote them originally
		result.sort(true)
		if err := compareIDLists(result, testCase.expected); err != nil {
			t.Errorf("Testing: %s. Got error: %v", testCase.testName, err)
		}
	}

}

func idFromInt(i int64) gouuidv6.UUID {
	return gouuidv6.NewFromTime(time.Unix(i, 0))
}

func compareIDLists(listA, listB idList) error {
	if len(listA) != len(listB) {
		return fmt.Errorf("Length of lists does not match. Got %v; wanted %v", len(listA), len(listB))
	}

	for i := range listA {
		if listA[i] != listB[i] {
			return fmt.Errorf("Comparing list members, mismatch at index %v. List A: %v; List B: %v", i, listA[i], listB[i])
		}
	}

	return nil
}
