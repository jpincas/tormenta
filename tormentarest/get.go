package tormentarest

import (
	"net/http"
	"reflect"

	"github.com/go-chi/chi"
	"github.com/jpincas/gouuidv6"
	tormenta "github.com/jpincas/tormenta/tormentadb"
	"github.com/jpincas/tormenta/utilities"
)

func getByID(w http.ResponseWriter, r *http.Request) {
	// Get the entity name from the URL,
	// look it up in the entity map,
	// then create a new one of that type to hold the results of the query
	idString := chi.URLParam(r, "id")
	entityName := entityNameFromPath_ID(r.URL.Path)
	entity := App.EntityMap[entityName]
	result := reflect.New(reflect.Indirect(reflect.ValueOf(entity)).Type()).Interface().(tormenta.Tormentable)

	id := gouuidv6.UUID{}
	if err := id.UnmarshalText([]byte(idString)); err != nil {
		renderError(w, utilities.ErrBadIDFormat, idString)
		return
	}

	// Get the record
	ok, err := App.DB.Get(result, id)
	if err != nil {
		renderError(w, utilities.ErrDBConnection)
		return
	}

	if !ok {
		renderError(w, utilities.ErrRecordNotFound, idString)
		return
	}

	// Render JSON
	App.Render.JSON(w, http.StatusOK, result)
}

func getList(w http.ResponseWriter, r *http.Request) {
	// Get the entity name from the URL,
	// look it up in the entity map,
	// then create a new slice of that type to hold the results of the query
	entityName := entityNameFromPath_List(r.URL.Path)
	entity := App.EntityMap[entityName]
	results := reflect.New(reflect.SliceOf(reflect.Indirect(reflect.ValueOf(entity)).Type())).Interface()

	// Set up the base query
	q := App.DB.Find(results)

	// Run the query builder,
	// to apply query options from the URL parameters
	if err := utilities.BuildQuery(r, q); err != nil {
		renderError(w, err.Error())
		return
	}

	// Run the query
	n, err := q.Run()
	if err != nil {
		renderError(w, utilities.ErrDBConnection)
		return
	}

	// Render JSON
	// For 0 results, render empty list
	if n == 0 {
		App.Render.JSON(w, http.StatusOK, []interface{}{})
		return
	}

	App.Render.JSON(w, http.StatusOK, results)
}
