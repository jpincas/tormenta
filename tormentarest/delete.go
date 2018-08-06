package tormentarest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta/utilities"
)

func deleteByID(w http.ResponseWriter, r *http.Request) {
	// Get the entity name from the URL,
	// look it up in the entity map,
	// then create a new one of that type to hold the results of the query
	idString := chi.URLParam(r, "id")
	entityName := entityNameFromPath_ID(r.URL.Path)

	id := gouuidv6.UUID{}
	if err := id.UnmarshalText([]byte(idString)); err != nil {
		renderError(w, utilities.ErrBadIDFormat, idString)
		return
	}

	// Delete the record
	n, err := App.DB.Delete(entityName, id)
	if err != nil {
		renderError(w, utilities.ErrDeleting, err)
		return
	}

	if n == 0 {
		renderError(w, utilities.ErrRecordNotFound, idString)
		return
	}

	// Render
	renderSuccess(w, n)
	return
}
