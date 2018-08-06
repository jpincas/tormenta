package utilities

//Error messages for GUI/REST layers
const (
	ErrDBConnection      = "Error connecting to DB"
	ErrSaving            = "Entity could not be saved. %v"
	ErrDeleting          = "Entity could not be deleted. %v"
	ErrGetting           = "Entity could not be retrieved. %v"
	ErrBadIDFormat       = "Bad format for Tormenta ID - %s"
	ErrRecordNotFound    = "Record with id %s not found"
	ErrBadLimitFormat    = "%s is an invalid input for LIMIT. Expecting a number"
	ErrBadOffsetFormat   = "%s is an invalid input for OFFSET. Expecting a number"
	ErrBadReverseFormat  = "%s is an invalid input for REVERSE. Expecting true/false"
	ErrBadFromFormat     = "Invalid input for FROM. Expecting somthing like '2006-01-02'"
	ErrBadToFormat       = "Invalid input for TO. Expecting somthing like '2006-01-02'"
	ErrIndexWithNoParams = "An index search has been specified, but no MATCH or START/END (for range) has been specified"
	ErrRangeTypeMismatch = "For a range index search, START and END should be of the same type (bool, int, float, string)"
	ErrUnmarshall        = "Error in format of data to save: %v"
)
