package tormenta_test

import (
	"fmt"
	"time"

	"github.com/jpincas/gouuidv6"
	"github.com/jpincas/tormenta"
)

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
	ProductID   gouuidv6.UUID
	Product     Product `tormenta:"-"`
}

func printlinef(formatString string, x interface{}) {
	fmt.Println(fmt.Sprintf(formatString, x))
}

func Example() {
	// Open the DB
	db, _ := tormenta.OpenTestWithOptions("data/tests", testDBOptions)
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
	printlinef("Saved %v records", n)

	// Get
	var nonExistentID gouuidv6.UUID
	var product Product

	// No such record
	ok, _ := db.Get(&product, nonExistentID)
	printlinef("Get record? %v", ok)

	// Get by entity
	ok, _ = db.Get(&product1)
	printlinef("Got record? %v", ok)

	// Get with optional separately specified ID
	ok, _ = db.Get(&product, product1.ID)
	printlinef("Get record with separately specified ID? %v", ok)

	// Delete
	db.Delete(&product1)
	fmt.Println("Deleted 1 record")

	// Basic query
	var products []Product
	n, _ = db.Find(&products).Run()
	printlinef("Found %v record(s)", n)

	// Date range query
	// Make some fullStructs with specific creation times
	var ttsToSave []tormenta.Record
	dates := []time.Time{
		// Specific years
		time.Date(2009, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2010, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2011, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2012, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2013, time.January, 1, 1, 0, 0, 0, time.UTC),
	}

	for i, date := range dates {
		ttsToSave = append(ttsToSave, &Order{
			// You wouln't normally do this manually
			// This is just for illustration
			Model: tormenta.Model{
				ID: gouuidv6.NewFromTime(date),
			},
			Customer:    fmt.Sprintf("customer-%v", i), // "customer-0", "customer-1"
			ShippingFee: float64(i),
		})
	}

	// Save the fullStructs
	db.Save(ttsToSave...)

	var fullStructs []Order
	var fullStruct Order

	mid2009 := time.Date(2009, time.June, 1, 1, 0, 0, 0, time.UTC)
	mid2010 := time.Date(2010, time.June, 1, 1, 0, 0, 0, time.UTC)
	mid2012 := time.Date(2012, time.June, 1, 1, 0, 0, 0, time.UTC)

	// Basic date range query
	n, _ = db.Find(&fullStructs).From(mid2009).To(mid2012).Run()
	printlinef("Basic date range query: %v records found", n)

	// First
	n, _ = db.First(&fullStruct).From(mid2009).To(mid2012).Run()
	printlinef("Basic date range query, first only: %v record(s) found", n)

	// First (not found)
	n, _ = db.First(&fullStruct).From(time.Now()).To(time.Now()).Run()
	printlinef("Basic date range query, first only: %v record(s) found", n)

	// Count only (fast!)
	c, _ := db.Find(&fullStructs).From(mid2009).To(mid2012).Count()
	printlinef("Basic date range query, count only: counted %v", c)

	// Limit
	n, _ = db.Find(&fullStructs).From(mid2009).To(mid2012).Limit(2).Run()
	printlinef("Basic date range query, 2 limit: %v record(s) found", n)

	// Offset
	n, _ = db.Find(&fullStructs).From(mid2009).To(mid2012).Limit(2).Offset(1).Run()
	printlinef("Basic date range query, 2 limit, 1 offset: %v record(s) found", n)

	// Reverse, count
	c, _ = db.Find(&fullStructs).Reverse().From(mid2009).To(mid2012).Count()
	printlinef("Basic date range query, reverse, count: %v record(s) counted", c)

	// Secondary index on 'customer' - exact index match
	n, _ = db.First(&fullStruct).Match("customer", "customer-2").Run()
	printlinef("Index query, exact match: %v record(s) found", n)

	// Secondary index on 'customer' - prefix match
	n, _ = db.First(&fullStruct).StartsWith("customer", "customer-").Run()
	printlinef("Index query, starts with: %v record(s) found", n)

	// Index range, QuickSum (based on index)
	var sum float64
	db.Find(&fullStructs).Range("shippingfee", 0.00, 10.00).From(mid2009).To(mid2012).QuickSum(&sum, "shippingfee")
	printlinef("Index range, date range, index sum query. Sum: %v", sum)

	// Secondary index on 'customer' - index range and count
	c, _ = db.Find(&fullStructs).Range("customer", "customer-1", "customer-3").Count()
	printlinef("Index range, count: %v record(s) counted", c)

	// Secondary index on 'customer' - exact index match, count and date range
	c, _ = db.Find(&fullStructs).Match("customer", "customer-3").From(mid2009).To(time.Now()).Count()
	printlinef("Index exact match, date range, count: %v record(s) counted", c)

	// Secondary index on 'customer' - index range AND date range
	c, _ = db.Find(&fullStructs).Range("customer", "customer-1", "customer-3").From(mid2009).To(mid2010).Count()
	printlinef("Index range, date range, count: %v record(s) counted", c)

	// Output:
	// Saved 2 records
	// Get record? false
	// Got record? true
	// Get record with separately specified ID? true
	// Deleted 1 record
	// Found 1 record(s)
	// Basic date range query: 3 records found
	// Basic date range query, first only: 1 record(s) found
	// Basic date range query, first only: 0 record(s) found
	// Basic date range query, count only: counted 3
	// Basic date range query, 2 limit: 2 record(s) found
	// Basic date range query, 2 limit, 1 offset: 2 record(s) found
	// Basic date range query, reverse, count: 3 record(s) counted
	// Index query, exact match: 1 record(s) found
	// Index query, starts with: 1 record(s) found
	// Index range, date range, index sum query. Sum: 6
	// Index range, count: 3 record(s) counted
	// Index exact match, date range, count: 1 record(s) counted
	// Index range, date range, count: 1 record(s) counted
}
