package tormentarest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// Serve serves a completely generic REST api over a Tormenta DB
func Serve(port string, db *tormenta.DB, entities ...tormenta.Tormentable) {
	// Initialise the application
	App.init(db)

	// Make the router
	r := New(entities...)

	// Show that we're starting
	fmt.Println(
		`
-----------------------------------
Starting Generic TormentaREST API...
-----------------------------------
		`)

	// Fire up the router
	http.ListenAndServe(port, r)
}

// ServeRouter serves a custom router over a Tormenta DB
func ServeRouter(r *chi.Mux, port string, db *tormenta.DB, entities ...tormenta.Tormentable) {
	// Initialise the application
	App.init(db)

	// Make the router
	WithRouter(r, entities...)

	// Show that we're starting
	fmt.Println(
		`
-----------------------------------
Starting Custom TormentaREST API...
-----------------------------------
		`)

	// Fire up the router and serve passed in router
	http.ListenAndServe(port, r)
}
