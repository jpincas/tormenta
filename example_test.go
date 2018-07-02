package tormenta

import (
	"log"
)

func Example_Main() {
	// Open the DB
	db, _ := OpenTest("data/tests")
	defer db.Close()

	// Create some products
	product1 := Product{
		Code:          "SKU1",
		Name:          "Product1",
		Price:         1.00,
		StartingStock: 50}
	product2 := Product{
		Code:          "SKU2",
		Name:          "Product2",
		Price:         2.00,
		StartingStock: 100}

	// Save them
	n, _ := db.Save(&product1, &product2)
	log.Println(n) // 2

	nonExistentID := newID()
	product1ID := product1.ID

	// Get by ID
	var product Product
	ok, _ := db.GetByID(&product, nonExistentID)
	log.Println(ok) // false
	ok, _ = db.GetByID(&product, product1ID)
	log.Println(ok) // true ( -> product)

	// Query
	var products []Product
	n, _ = db.Query(&products).Run()
	log.Println(n) // 2 (-> products)

}
