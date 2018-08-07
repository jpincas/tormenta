package tormentarest

import (
	"github.com/go-chi/chi"
	tormenta "github.com/jpincas/tormenta/tormentadb"
	"github.com/jpincas/tormenta/utilities"
)

// WithRouter builds a REST api on the specified router
func WithRouter(r *chi.Mux, db *tormenta.DB, root string, entities ...tormenta.Tormentable) {
	// Initialise the application
	App.init(db, root)
	// Build the router
	buildRouter(r, root, entities...)
}

func buildRouter(r *chi.Mux, root string, entities ...tormenta.Tormentable) {

	r.Route("/"+root, func(r chi.Router) {

		for _, entity := range entities {
			entityName := utilities.Pluralise(tormenta.KeyRootString(entity))
			App.EntityMap[entityName] = entity

			r.Route("/"+entityName, func(r chi.Router) {
				// GET
				r.Get("/", getList)
				r.Get("/{id}", getByID)

				// DELETE
				r.Delete("/{id}", deleteByID)

				// POST
				r.Post("/", post)

				// PUT
				r.Put("/{id}", putByID)
			})
		}
	})

}
