package tormentagui

import (
	"net/http"
	"reflect"
	"strings"

	tormenta "github.com/jpincas/tormenta/tormentadb"
	"github.com/jpincas/tormenta/utilities"
)

// func getByID(w http.ResponseWriter, r *http.Request) {
// 	// Get the entity name from the URL,
// 	// look it up in the entity map,
// 	// then create a new one of that type to hold the results of the query
// 	idString := chi.URLParam(r, "id")
// 	entityName := strings.Split(r.URL.Path, "/")[1]
// 	entity := App.EntityMap[entityName]
// 	result := reflect.New(reflect.Indirect(reflect.ValueOf(entity)).Type()).Interface().(tormenta.Tormentable)

// 	id := gouuidv6.UUID{}
// 	if err := id.UnmarshalText([]byte(idString)); err != nil {
// 		//
// 		return
// 	}

// 	// Get the record
// 	ok, err := App.DB.Get(result, id)
// 	if err != nil {
// 		//
// 		return
// 	}

// 	if !ok {
// 		//
// 		return
// 	}

// 	// Render JSON
// 	App.Render.JSON(w, http.StatusOK, result)
// }

func listEntities(w http.ResponseWriter, r *http.Request) {
	templateData := struct {
		Entities map[string]tormenta.Tormentable
	}{
		Entities: App.EntityMap,
	}

	App.Templates.ExecuteTemplate(w, "entities.html", templateData)
}

func getList(w http.ResponseWriter, r *http.Request) {
	// Get the entity name from the URL,
	// look it up in the entity map,
	// then create a new slice of that type to hold the results of the query
	entityName := strings.TrimPrefix(r.URL.Path, "/")
	entity := App.EntityMap[entityName]
	results := reflect.New(reflect.SliceOf(reflect.Indirect(reflect.ValueOf(entity)).Type())).Interface()

	// Set up the base query
	q := App.DB.Find(results)

	// Run the query builder,
	// to apply query options from the URL parameters
	if err := utilities.BuildQuery(r, q); err != nil {
		App.Templates.ExecuteTemplate(w, "error-partial.html", struct{}{})
		return
	}

	// Run the query
	n, err := q.Run()
	if err != nil {
		App.Templates.ExecuteTemplate(w, "error-partial.html", struct{}{})
		return
	}

	// Render JSON
	// For 0 results, render empty list
	if n == 0 {
		App.Templates.ExecuteTemplate(w, "error-partial.html", struct{}{})
		return
	}

	templateData := struct {
		Results interface{}
		Fields  []string
	}{
		Results: results,
		Fields:  tormenta.ListFields(entity),
	}

	App.Templates.ExecuteTemplate(w, "list.html", templateData)
}
