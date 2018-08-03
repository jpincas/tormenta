package tormentagui

import (
	"fmt"
	"html/template"
	"os"

	tormenta "github.com/jpincas/tormenta/tormentadb"
)

const packagePath = "github.com/jpincas/tormenta/tormentagui"

// App is the application wide construct
var App application

type application struct {
	DB        *tormenta.DB
	Templates *template.Template
	EntityMap map[string]tormenta.Tormentable
	GoPath    string
}

func (a application) templatesLocation() string {
	return fmt.Sprintf("%s/src/%s/templates/*.html", a.GoPath, packagePath)
}

func (a application) staticFilesLocation() string {
	return fmt.Sprintf("%s/src/%s/static", a.GoPath, packagePath)
}

func (a *application) parseTemplates() {
	//Formatting functions for templates
	funcMap := template.FuncMap{
		"printJSONWithoutModel": printJSONWithoutModel,
		"mapFields":             tormenta.MapFields,
		"autoFormat":            autoFormat,
	}

	a.Templates = template.Must(template.New("main").Funcs(funcMap).ParseGlob(a.templatesLocation()))
}

func (a *application) init(db *tormenta.DB) {
	a.DB = db
	a.EntityMap = map[string]tormenta.Tormentable{}
	a.GoPath = os.Getenv("GOPATH")

	a.parseTemplates()
}
