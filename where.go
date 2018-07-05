package tormenta

type whereFilterFunction = func(s string) bool

type Filter struct {
	index    string
	function whereFilterFunction
}

// func applyWhereFilters(key []byte, filter Filter) bool {
// 	for _, filter := range filters {
// 		passes := filter.function(filter.index)
// 	}

// 	return true
// }
