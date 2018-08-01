package tormentarest

import (
	tormenta "github.com/jpincas/tormenta/tormentadb"
	"github.com/unrolled/render"
)

// App is the application wide construct
var App application

type application struct {
	DB        *tormenta.DB
	Render    *render.Render
	EntityMap map[string]tormenta.Tormentable
}

func (a *application) init(db *tormenta.DB) {
	// DB
	a.DB = db

	// Renderer
	a.Render = render.New()

	// Entity Map
	a.EntityMap = map[string]tormenta.Tormentable{}
}

// Init initialises TormentaREST to the passed in DB
func Init(db *tormenta.DB) {
	App.init(db)
}
