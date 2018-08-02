package tormentagui

import (
	"fmt"
	"net/http"

	"github.com/jpincas/tormenta/tormentarest"

	"github.com/go-chi/chi"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// Serve serves a completely generic REST api over a Tormenta DB
func Serve(port string, db *tormenta.DB, entities ...tormenta.Tormentable) {
	// Initialise the application
	App.init(db)

	// Make the router
	r := chi.NewRouter()
	buildRouter(r, db, entities...)

	// Show that we're starting
	fmt.Println(
		`
-----------------------------------
	    Starting Tormenta GUI
-----------------------------------
		`)

	// Fire up the router
	http.ListenAndServe(port, r)
}

func buildRouter(r *chi.Mux, db *tormenta.DB, entities ...tormenta.Tormentable) {

	r.Get("/", home)

	// Include Tormenta REST
	tormentarest.WithRouter(r, db, "api", entities...)

	for _, entity := range entities {
		entityName := tormenta.KeyRootString(entity)
		App.EntityMap[entityName] = entity

		r.Route("/"+entityName, func(r chi.Router) {
			// GET
			r.Get("/", getList)
			r.Get("/{id}", getByID)
		})
	}

	// Static file server
	fileServer(r, "/static", http.Dir(App.staticFilesLocation()))
}
