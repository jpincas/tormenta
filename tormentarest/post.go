package tormentarest

import (
	"encoding/json"
	"net/http"
	"reflect"

	tormenta "github.com/jpincas/tormenta/tormentadb"
	"github.com/jpincas/tormenta/utilities"
)

func post(w http.ResponseWriter, r *http.Request) {
	// Get the entity name from the URL,
	// look it up in the entity map,
	// then create a new one of that type to hold the new entity
	entityName := entityNameFromPath_ID(r.URL.Path)
	entity := App.EntityMap[entityName]
	toSave := reflect.New(reflect.Indirect(reflect.ValueOf(entity)).Type()).Interface().(tormenta.Tormentable)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&toSave)
	if err != nil {
		renderError(w, utilities.ErrUnmarshall)
		return
	}

	_, err = App.DB.Save(toSave)
	if err != nil {
		renderError(w, utilities.ErrDBConnection)
		return
	}

	App.Render.JSON(w, http.StatusOK, toSave)
	return
}
