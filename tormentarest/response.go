package tormentarest

import (
	"fmt"
	"net/http"
)

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}
type noProcessedResponse struct {
	NoRecords int `json:"noRecords"`
}

func renderError(w http.ResponseWriter, s string, i ...interface{}) {
	App.Render.JSON(w, http.StatusBadRequest, errorResponse{
		ErrorMessage: fmt.Sprintf(s, i...),
	})
}

func renderSuccess(w http.ResponseWriter, n int) {
	App.Render.JSON(w, http.StatusOK, noProcessedResponse{
		NoRecords: n,
	})
}
