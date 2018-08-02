package main

import (
	"log"

	"github.com/go-chi/chi"
	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
	"github.com/jpincas/tormenta/tormentarest"
)

func main() {
	serveGeneric()
	// serveCustom()
}

func serveGeneric() {
	// Open the DB
	db, err := tormenta.OpenTest("data")
	if err != nil {
		log.Fatalf("Could not open DB: %s", err)
		return
	}

	// Populate with data
	demo.PopulateWithDemoData(db, 100)

	// Serve the REST api
	tormentarest.Serve(
		"api",   // root
		":3333", // the port you want to serve the api on
		db,      // connection to the Tormenta DB
		// List of entities to include in the API
		&demo.Order{},   // -> /order
		&demo.Product{}, // -> /product
	)
}

func serveCustom() {
	// Open the DB
	db, err := tormenta.OpenTest("data")
	if err != nil {
		log.Fatalf("Could not open DB: %s", err)
		return
	}

	// Populate with data
	demo.PopulateWithDemoData(db, 100)

	// New router
	r := chi.NewRouter()

	// Serve the REST api
	tormentarest.ServeRouter(r, "api", ":3333", db, &demo.Order{}, &demo.Product{})
}
