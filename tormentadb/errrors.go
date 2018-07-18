package tormentadb

// Error messages
const (
	ErrNilInputMatchIndexQuery        = "Nil is not a valid input for an exact match search"
	ErrNilInputsRangeIndexQuery       = "Nil from both ends of the range is not a valid input for an index range search"
	ErrMoreThan2InputsRangeIndexQuery = "Index Query Where clause requires either 1 (exact match) or 2 (range) parameters"
)
