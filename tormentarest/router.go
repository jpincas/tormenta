package tormentarest

import (
	"github.com/go-chi/chi"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// New takes a list of entities and constucts the REST endpoints
func New(entities ...tormenta.Tormentable) *chi.Mux {
	r := chi.NewRouter()
	buildRouter(r, entities...)
	return r
}

// WithRouter builds a REST api on the specified router
func WithRouter(r *chi.Mux, entities ...tormenta.Tormentable) {
	buildRouter(r, entities...)
}

func buildRouter(r *chi.Mux, entities ...tormenta.Tormentable) {
	for _, entity := range entities {
		entityName := entityRoot(entity)
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
}
