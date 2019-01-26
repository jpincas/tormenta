package example

import (
	"fmt"
	"log"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
)

// Include 'go:generate msgp' in your file and run 'go generate' to generate MessagePack marshall/unmarshall methods

// Define your data.
// Include tormenta.Model to get date ordered IDs, last updated field etc
// Tag with 'tormenta:"noindex"' to skip secondary index creation
type Product struct {
	tormenta.Model
	Code          string
	Name          string
	Price         float32
	StartingStock int
	Tags          []string
}

type Order struct {
	tormenta.Model
	Customer    string
	Department  int
	ShippingFee float64
}

func Example_Main() {
	// Open the DB
	db, _ := tormenta.OpenTest("data/tests")
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

	// Save
	n, _ := db.Save(&product1, &product2)
	log.Println("Saved: ", n) // 2

	// Get
	var nonExistentID gouuidv6.UUID
	var product Product

	// No such record
	ok, _ := db.Get(&product, nonExistentID)
	log.Println("Get: ", ok) // false

	// Get by entity ID
	ok, _ = db.Get(&product1)
	log.Println("Get entity: ", ok) // true ( -> product 1)

	// Get with optional separately specified ID
	ok, _ = db.Get(&product, product1.ID)
	log.Println("Get by entity ID: ", ok) // true ( -> product 1)

	// Delete
	n, _ = db.Delete("product", product1.ID)
	log.Println("Delete: ", n) // 1

	// Basic query
	var products []Product
	n, _ = db.Find(&products).Run()
	log.Println("Find: ", n) // 2 (-> products)

	// Date range query
	// Make some orders with specific creation times
	var ordersToSave []tormenta.Tormentable
	dates := []time.Time{
		// Specific years
		time.Date(2009, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2010, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2011, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2012, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2013, time.January, 1, 1, 0, 0, 0, time.UTC),
	}

	for i, date := range dates {
		ordersToSave = append(ordersToSave, &Order{
			// You wouln't normally do this manually
			// This is just for illustration
			Model: tormenta.Model{
				ID: gouuidv6.NewFromTime(date),
			},
			Customer:    fmt.Sprintf("customer-%v", i), // "customer-0", "customer-1"
			ShippingFee: float64(i),
		})
	}

	// Save the orders
	db.Save(ordersToSave...)

	var orders []Order
	var order Order

	mid2009 := time.Date(2009, time.June, 1, 1, 0, 0, 0, time.UTC)
	mid2010 := time.Date(2010, time.June, 1, 1, 0, 0, 0, time.UTC)
	mid2012 := time.Date(2012, time.June, 1, 1, 0, 0, 0, time.UTC)

	// Basic date range query
	n, _ = db.Find(&orders).From(mid2009).To(mid2012).Run()
	log.Println("Basic - date range: ", n) // 3 (-> orders )

	// First
	n, _ = db.First(&order).From(mid2009).To(mid2012).Run()
	log.Println("First - found: ", n) // 1 (-> order )

	// First (not found)
	n, _ = db.First(&order).From(time.Now()).To(time.Now()).Run()
	log.Println("First - not found: ", n) // 0

	// Count only (fast!)
	c, _ := db.Find(&orders).From(mid2009).To(mid2012).Count()
	log.Println("Count: ", c) // 3

	// Limit
	n, _ = db.Find(&orders).From(mid2009).To(mid2012).Limit(2).Run()
	log.Println("Limit: ", n) // 2

	// Offset
	n, _ = db.Find(&orders).From(mid2009).To(mid2012).Limit(2).Offset(1).Run()
	log.Println("Limit and offset: ", n) // 2 (same count, different results)

	// Reverse
	c, _ = db.Find(&orders).Reverse().From(mid2009).To(mid2012).Count()
	log.Println("Reverse: ", c) // 3

	// Secondary index on 'customer' - exact index match
	n, _ = db.First(&order).Match("customer", "customer-2").Run()
	log.Println("Index - exact match: ", n) // 1 (-> order )

	// Secondary index on 'customer' - prefix match
	n, _ = db.First(&order).StartsWith("customer", "customer-").Run()
	log.Println("Index - prefix match: ", n) // 5 (-> order )

	// Sum (based on index)
	var sum float64
	db.Find(&orders).Range("shippingfee", 0.00, 10.00).From(mid2009).To(mid2012).Sum(&sum)
	log.Println("Sum: ", sum) // 6.00 (1.00 + 2.00 + 3.00)

	// Secondary index on 'customer' - index range and count
	c, _ = db.Find(&orders).Range("customer", "customer-1", "customer-3").Count()
	log.Println("Index - range: ", c) // 3

	// Secondary index on 'customer' - exact index match, count and date range
	c, _ = db.Find(&orders).Match("customer", "customer-3").From(mid2009).To(time.Now()).Count()
	log.Println("Index - exact match and date range: ", c) // 1

	// Secondary index on 'customer' - index range AND date range
	c, _ = db.Find(&orders).Range("customer", "customer-1", "customer-3").From(mid2009).To(mid2010).Count()
	log.Println("Index - range and date range", c) // 1
}
