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

	// Try unmarhsalling to uuuidv6
	// But don't worry if you can't, we're actually going to use that case
	// for something useful
	id := gouuidv6.UUID{}
	id.UnmarshalText([]byte(idString))

	// Get the record
	// If you can't find it, again, don't worry
	_, timer, err := App.DB.Get(result, id)
	if err != nil {
		App.Templates.ExecuteTemplate(w, "error-partial.html", struct{}{})
		return
	}

	// At this point, the result will be correct
	// if a correct ID has been provided
	// Otherwise, the zero-value (blank struct) will be there
	// and this will get marhsalled and return
	// We can use this as a tempalte for creating new entities!
	// i.e. if I call /api/order/new
	// 'new' doesnt exist, so I get the correct blank template back
	templateData := struct {
		Result     interface{}
		EntityName string
		Timer      int
	}{
		Result:     result,
		EntityName: entityName,
		Timer:      timer,
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
	// Reverse as default
	q := App.DB.Find(results)

	// IF this is NOT a query, we will reverse it by default,
	// so that the most recent records are at the top
	// IF this is a query (where the 'reverse' param has been set as true or false)
	// then we will respect that setting
	if !utilities.HasReverseBeenSet(r) {
		q.Reverse()
	}

	// Run the query builder,
	// to apply query options from the URL parameters
	if err := utilities.BuildQuery(r, q); err != nil {
		App.Templates.ExecuteTemplate(w, "error-partial.html", struct{}{})
		return
	}

	// Run the query
	n, timer, err := q.Run()
	if err != nil {
		App.Templates.ExecuteTemplate(w, "error-partial.html", struct{}{})
		return
	}

	// Template override with URL param
	template := "results"
	overrideTemplate := r.URL.Query().Get("template")
	if r.URL.Query().Get("template") != "" {
		template = overrideTemplate
	}

	templateData := struct {
		NoResults  int
		Results    interface{}
		Fields     []string
		EntityName string
		Timer      int
	}{
		NoResults:  n,
		Results:    results,
		Fields:     tormenta.ListFields(entity),
		EntityName: entityName,
		Timer:      timer,
	}

	App.Templates.ExecuteTemplate(w, template+".html", templateData)
}
