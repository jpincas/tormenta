package tormentarest

import (
	"fmt"
	"net/http"
)

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

const (
	errDBConnection     = "Error connecting to DB"
	errBadIDFormat      = "Bad format for Tormenta ID - %s"
	errRecordNotFound   = "Record with id %s not found"
	errBadLimitFormat   = "%s is an invalid input for LIMIT. Expecting a number"
	errBadReverseFormat = "%s is an invalid input for REVERSE. Expecting true/false"
	errBadFromFormat    = "Invalid input for FROM. Expecting somthing like '2006-01-02'"
	errBadToFormat      = "Invalid input for TO. Expecting somthing like '2006-01-02'"
)

func renderError(w http.ResponseWriter, s string, i ...interface{}) {
	App.Render.JSON(w, http.StatusBadRequest, errorResponse{
		ErrorMessage: fmt.Sprintf(s, i...),
	})
}
