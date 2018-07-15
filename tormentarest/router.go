package tormentadbrest

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jpincas/tormenta"
)

// func main() {
// 	r := chi.NewRouter()

// 	// RESTy routes for "articles" resource
// 	r.Route("/{entityType}", func(r chi.Router) {
// 		r.Get("/", listAll) // GET /articles
// 		// r.Post("/", createNew) // POST /articles

// 		// r.Route("/{entityID}", func(r chi.Router) {
// 		// 	r.Get("/", getByID)       // GET /articles/123
// 		// 	r.Put("/", updateByID)    // PUT /articles/123
// 		// 	r.Delete("/", deleteByID) // DELETE /articles/123
// 		// })
// 	})

// 	http.ListenAndServe(":3333", r)
// }

func makeRouter(entities ...tormenta.Tormentable) *chi.Mux {
	r := chi.NewRouter()

	for range entities {

		r.Route("/test", func(r chi.Router) {
			r.Get("/", listAll) // GET /articles
			// r.Post("/", createNew) // POST /articles

			// r.Route("/{entityID}", func(r chi.Router) {
			// 	r.Get("/", getByID)       // GET /articles/123
			// 	r.Put("/", updateByID)    // PUT /articles/123
			// 	r.Delete("/", deleteByID) // DELETE /articles/123
			// })
		})

	}

	return r
}

func listAll(w http.ResponseWriter, r *http.Request) {

}
