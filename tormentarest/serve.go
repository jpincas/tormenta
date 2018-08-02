package tormentarest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// Serve serves a completely generic REST api over a Tormenta DB
func Serve(root, port string, db *tormenta.DB, entities ...tormenta.Tormentable) {
	r := chi.NewRouter()
	ServeRouter(r, root, port, db, entities...)
}

// ServeRouter serves a custom router over a Tormenta DB
func ServeRouter(r *chi.Mux, root, port string, db *tormenta.DB, entities ...tormenta.Tormentable) {

	// Make the router
	WithRouter(r, db, root, entities...)

	// Show that we're starting
	fmt.Println(
		`
-----------------------------------
Starting TormentaREST API...
-----------------------------------
		`)

	// Fire up the router and serve passed in router
	err := http.ListenAndServe(port, r)
	if err != nil {
		panic("Could not start TormentaREST server")
	}
}
