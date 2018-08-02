package tormentagui

import (
	"encoding/json"
	"fmt"
	"html/template"

	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// App is the application wide construct
var App application

type application struct {
	DB        *tormenta.DB
	Templates *template.Template
	EntityMap map[string]tormenta.Tormentable
}

// prettyPrintStruct just outputs a struct as indented JSON
func prettyPrintStruct(s interface{}) string {
	res, _ := json.MarshalIndent(s, "", " ")
	return fmt.Sprintln(string(res))
}

func (a *application) parseTemplates() {
	//Formatting functions for templates
	funcMap := template.FuncMap{
		"prettyPrint": prettyPrintStruct,
		"mapFields":   tormenta.MapFields,
		"autoFormat":  autoFormat,
	}

	a.Templates = template.Must(template.New("main").Funcs(funcMap).ParseGlob("../templates/*.html"))
}

func (a *application) init(db *tormenta.DB) {
	a.DB = db
	a.EntityMap = map[string]tormenta.Tormentable{}

	a.parseTemplates()
}
