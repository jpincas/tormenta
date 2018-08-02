package main

import (
	"log"

	"github.com/jpincas/tormenta/demo"
	tormenta "github.com/jpincas/tormenta/tormentadb"
	"github.com/jpincas/tormenta/tormentagui"
)

func main() {
	// Open the DB
	db, err := tormenta.OpenTest("data")
	if err != nil {
		log.Fatalf("Could not open DB: %s", err)
		return
	}

	// Populate with data
	demo.PopulateWithDemoData(db, 100)

	// Serve the REST api
	tormentagui.Serve(
		":3334", // the port you want to serve the GUI on
		db,      // connection to the Tormenta DB
		// List of entities to include in the API
		&demo.Order{},   // -> /order
		&demo.Product{}, // -> /product
	)
}
