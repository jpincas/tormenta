package tormenta

// Error messages
const (
	ErrNilInputMatchIndexQuery   = "Nil is not a valid input for an exact match search"
	ErrNilInputsRangeIndexQuery  = "Nil from both ends of the range is not a valid input for an index range search"
	ErrBlankInputStartsWithQuery = "Blank string is not valid input for 'starts with' query"
	ErrFieldCouldNotBeFound      = "Field %s could not be found"
)
