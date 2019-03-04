package tormenta

import (
	"sort"

	"github.com/jpincas/gouuidv6"
)

type idList []gouuidv6.UUID

func (ids idList) length() int {
	return len(ids)
}

func (ids idList) sort(reverse bool) {
	compareFunc := func(i, j int) bool {
		return ids[i].Compare(ids[j])
	}

	if reverse {
		compareFunc = func(i, j int) bool {
			return ids[j].Compare(ids[i])
		}
	}

	sort.Slice(ids, compareFunc)
}

// for OR
func union(listsOfIDs ...idList) (result idList) {
	masterMap := map[gouuidv6.UUID]bool{}

	for _, list := range listsOfIDs {
		for _, id := range list {
			masterMap[id] = true
		}
	}

	for id := range masterMap {
		result = append(result, id)
	}

	return result
}

// for AND
func intersection(listsOfIDs ...idList) (result idList) {
	// Deal with emtpy and single list cases
	if len(listsOfIDs) == 0 {
		return
	}

	if len(listsOfIDs) == 1 {
		result = listsOfIDs[0]
		return
	}

	// Map out the IDs from each list,
	// keeping a count of how many times each has appeared in a list
	// In order that duplicates within a list don't count twice, we use a nested
	// map to keep track of the contributions from the currently iterating list
	// and only accept each IDs once
	masterMap := map[gouuidv6.UUID]int{}
	for _, list := range listsOfIDs {

		thisListIDs := map[gouuidv6.UUID]bool{}

		for _, id := range list {
			if _, found := thisListIDs[id]; !found {
				thisListIDs[id] = true
				masterMap[id] = masterMap[id] + 1
			}
		}
	}

	// Only append an ID to the list if it has appeared
	// in all the lists
	for id, count := range masterMap {
		if count == len(listsOfIDs) {
			result = append(result, id)
		}
	}

	return
}

var (
	fixedID1 = gouuidv6.New()
	fixedID2 = gouuidv6.New()
)

// isOr is a simple function to tell you wether a given combinator function is a union (or) or not
func isOr(combinator func(...idList) idList) bool {
	n := combinator([]idList{
		idList{fixedID1},
		idList{fixedID2},
	}...)

	// If the combinator is AND, then it will want the ID to appear in both
	// lists, thus giving a length of 0
	return len(n) != 0
}
