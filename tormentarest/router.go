package tormentarest

import (
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi"
	tormenta "github.com/jpincas/tormenta/tormentadb"
	"github.com/unrolled/render"
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
var App application

type application struct {
	DB        *tormenta.DB
	Render    *render.Render
	EntityMap map[string]tormenta.Tormentable
}

func MakeRouter(entities ...tormenta.Tormentable) *chi.Mux {
	r := chi.NewRouter()

	for _, entity := range entities {
		entityName := entityRoot(entity)
		App.EntityMap[entityName] = entity

		r.Route("/"+entityName, func(r chi.Router) {
			r.Get("/", GetList)
		})
	}

	// r.Route("/order", func(r chi.Router) {

	// 	r.Get("/order", listAll)
	// 	// r.Post("/", createNew) // POST /articles

	// 	// r.Route("/{entityID}", func(r chi.Router) {
	// 	// 	r.Get("/", getByID)       // GET /articles/123
	// 	// 	r.Put("/", updateByID)    // PUT /articles/123
	// 	// 	r.Delete("/", deleteByID) // DELETE /articles/123
	// 	// })
	// })

	return r
}

func (a *application) init(db *tormenta.DB) {
	// DB
	a.DB = db

	// Renderer
	a.Render = render.New()

	// Entity Map
	a.EntityMap = map[string]tormenta.Tormentable{}
}

// Serve serves a generic REST api over a Tormenta DB
func Serve(port string, db *tormenta.DB, entities ...tormenta.Tormentable) {
	// Initialise the application
	App.init(db)

	// Make the router
	r := MakeRouter(entities...)

	// Show that we're starting
	log.Println("Starting TormentaREST...")

	// Fire up the router
	http.ListenAndServe(port, r)
}

func entityRoot(entity tormenta.Tormentable) string {
	return string(tormenta.KeyRoot(entity))
}

func GetList(w http.ResponseWriter, r *http.Request) {
	// Get the entity name from the URL,
	// look it up in the entity map,
	// then create a new slice of that type to hold the results of the query
	entityName := strings.TrimPrefix(r.URL.Path, "/")
	entity := App.EntityMap[entityName]
	results := reflect.New(reflect.SliceOf(reflect.Indirect(reflect.ValueOf(entity)).Type())).Interface()

	// Run the query
	App.DB.Find(results).Run()

	// Render JSON
	App.Render.JSON(w, http.StatusOK, results)
}
