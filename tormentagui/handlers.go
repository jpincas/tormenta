package tormentagui

import (
	"net/http"
	"reflect"
	"strings"

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
	entityName := strings.Split(r.URL.Path, "/")[1]
	entity := App.EntityMap[entityName]
	result := reflect.New(reflect.Indirect(reflect.ValueOf(entity)).Type()).Interface().(tormenta.Tormentable)

	id := gouuidv6.UUID{}
	if err := id.UnmarshalText([]byte(idString)); err != nil {
		//
		return
	}

	// Get the record
	ok, err := App.DB.Get(result, id)
	if err != nil {
		//
		return
	}

	if !ok {
		//
		return
	}

	templateData := struct {
		Result     interface{}
		EntityName string
	}{
		Result:     result,
		EntityName: entityName,
	}

	App.Templates.ExecuteTemplate(w, "detail.html", templateData)
}

func home(w http.ResponseWriter, r *http.Request) {
	templateData := struct {
		Entities map[string]tormenta.Tormentable
	}{
		Entities: App.EntityMap,
	}

	App.Templates.ExecuteTemplate(w, "home.html", templateData)
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

	templateData := struct {
		NoResults  int
		Results    interface{}
		Fields     []string
		EntityName string
	}{
		NoResults:  n,
		Results:    results,
		Fields:     tormenta.ListFields(entity),
		EntityName: entityName,
	}

	App.Templates.ExecuteTemplate(w, "list.html", templateData)
}
