package tormentarest

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi"
	"github.com/jpincas/gouuidv6"
	tormenta "github.com/jpincas/tormenta/tormentadb"
)

func putByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	entityName := strings.Split(r.URL.Path, "/")[1]
	entity := App.EntityMap[entityName]
	toSave := reflect.New(reflect.Indirect(reflect.ValueOf(entity)).Type()).Interface().(tormenta.Tormentable)

	id := gouuidv6.UUID{}
	if err := id.UnmarshalText([]byte(idString)); err != nil {
		renderError(w, errBadIDFormat, idString)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&toSave)
	if err != nil {
		renderError(w, errUnmarshall)
		return
	}

	// Once we have decoded the body,
	// set the ID to the ID specified in the URL
	// The upshot of this is that any ID specified in the body will be overwritten

	// Reminder: if the record exists, the contents will be overwritten
	// Otherwise it will be created
	toSave.SetID(id)

	_, err = App.DB.Save(toSave)
	if err != nil {
		renderError(w, errDBConnection)
		return
	}

	App.Render.JSON(w, http.StatusOK, toSave)
	return
}
