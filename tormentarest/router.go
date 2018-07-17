package tormentarest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// MakeRouter takes a list of entities and constucts the REST endpoints
func MakeRouter(entities ...tormenta.Tormentable) *chi.Mux {
	r := chi.NewRouter()

	for _, entity := range entities {
		entityName := entityRoot(entity)
		App.EntityMap[entityName] = entity

		r.Route("/"+entityName, func(r chi.Router) {
			r.Get("/", getList)
			r.Get("/{id}", getByID)
		})
	}

	return r
}

// Serve serves a generic REST api over a Tormenta DB
func Serve(port string, db *tormenta.DB, entities ...tormenta.Tormentable) {
	// Initialise the application
	App.init(db)

	// Make the router
	r := MakeRouter(entities...)

	// Show that we're starting
	fmt.Println(
		`
------------------------
Starting TormentaREST...
------------------------
		`)

	// Fire up the router
	http.ListenAndServe(port, r)
}
