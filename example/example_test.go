package example

import (
	"fmt"
	"log"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
)

// go:generate msgp
// Include 'go:generate msgp' in your file and run 'go generate' to generate MessagePack marshall/unmarshall methods

// Define your data.
// Include tormenta.Model to get date ordered IDs, last updated field etc
// Tag with 'tormenta:"index"' to create secondary indexes
type Product struct {
	tormenta.Model
	Code          string
	Name          string `tormenta:"index"`
	Price         float32
	StartingStock int
}

type Order struct {
	tormenta.Model
	Customer    string  `tormenta:"index"`
	Department  int     `tormenta:"index"`
	ShippingFee float64 `tormenta:"index"`
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

	// Save them
	n, _ := db.Save(&product1, &product2)
	log.Println(n) // 2

	// Get by ID
	var nonExistentID gouuidv6.UUID
	product1ID := product1.ID

	var product Product
	ok, _ := db.GetByID(&product, nonExistentID)
	log.Println(ok) // false

	ok, _ = db.GetByID(&product, product1ID)
	log.Println(ok) // true ( -> product)

	// Basic query
	var products []Product
	n, _ = db.Find(&products).Run()
	log.Println(n) // 2 (-> products)

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
			Customer: fmt.Sprintf("customer-%v", i), // "customer-0", "customer-1"
		})
	}

	// Save the orders
	db.Save(ordersToSave...)

	var orders []Order
	var order Order

	mid2009 := time.Date(2009, time.June, 1, 1, 0, 0, 0, time.UTC)
	mid2012 := time.Date(2012, time.June, 1, 1, 0, 0, 0, time.UTC)

	// Basic date range query
	n, _ = db.Find(&orders).From(mid2009).To(mid2012).Run()
	log.Println(n) // 3 (-> orders )

	// First
	n, _ = db.First(&order).From(mid2009).To(mid2012).Run()
	log.Println(n) // 1 (-> order )

	// First (not found)
	n, _ = db.First(&order).From(time.Now()).To(time.Now()).Run()
	log.Println(n) // 0

	// Count only (fast!)
	count, _ := db.Find(&orders).From(mid2009).To(mid2012).Count()
	log.Println(count) // 3

	// Limit
	n, _ = db.Find(&orders).From(mid2009).To(mid2012).Limit(2).Run()
	log.Println(n) // 2

	// Offset
	n, _ = db.Find(&orders).From(mid2009).To(mid2012).Limit(2).Offset(1).Run()
	log.Println(n) // 2 (same count, different results)

	// Reverse (note reversed range)
	count, _ = db.Find(&orders).Reverse().From(mid2012).To(mid2009).Count()
	log.Println(count) // 3

	// Secondary index on 'customer' - exact index match
	n, _ = db.First(&order).Where("customer", "customer-2").Run()
	log.Println(n) // 1 (-> order )

	// Secondary index on 'customer' - index range and count
	count, _ = db.Find(&orders).Where("customer", "customer-1", "customer-3").Count()
	log.Println(count) // 3

	// Secondary index on 'customer' - exact index match, count and date range
	count, _ = db.Find(&orders).Where("customer", "customer-3").From(mid2009).To(time.Now()).Count()
	log.Println(count) // 3

	// Secondary index on 'customer' - index range AND date range
	// Coming soon!
}
