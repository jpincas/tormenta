# TormentaDB (WIP)

Tormenta is a thin functionality layer over BadgerDB key/value store.  It provides simple, embedded object persistence for Go projects with some limited data querying capabilities and ORM-like features.  It uses date-ordered IDs so is particuarly good for data sets that are natually chronological, like financial transactions, soical media posts etc. Powered by:

- [BadgerDB](https://github.com/dgraph-io/badger)
- [TinyLib MessagePack](https://github.com/tinylib/msgp)
- ['V6' UUIDs](https://github.com/bradleypeabody/gouuidv6)
 
and greatly inspired by [Storm](https://github.com/asdine/storm).

## Quick How To

Add `tormenta.Model` to structs you want to persist and `tormenta:"index"` to fields you want to create secondary indexes for, install [TinyLib MessagePack codegen tool](https://github.com/tinylib/msgp), add `//go:generate msgp` to the top of your type definition files and run `go generate` whenever you change your structs.

Open a DB connection with `db, err := tormenta.Open("mydatadirectory")` (dont forget to `defer db.Close()`).

Save a single entity with `db.Save(&MyEntity)` or multiple entities with `db.Save(&MyEntity1, &MyEntity2)`.

Get a single entity by ID with `db.GetByID(&MyEntity, entityID)`.

Construct a query to find single or mutliple entities with `db.First(&MyEntity)` or `db.Find(&MyEntities)` respectively. Build up the query by chaining methods: `.From()/.To()` to add a date range, `.Where("indexName", indexStartOrExactMatch, optionalIndexEnd)` to add an index clause (exact match or range), `.Reverse()` to reverse the order of searching/results and `.Limit()/.Offset()` to limit the number of results.

Kick off the query with `.Run()`, or `.Count()` if you just need the count.  `.Sum()` is also available for float/int index searches.
	
## Example

```go
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
	log.Println("Saved: ", n) // 2

	// Get by ID
	var nonExistentID gouuidv6.UUID
	product1ID := product1.ID

	var product Product
	ok, _ := db.GetByID(&product, nonExistentID)
	log.Println("Get: ", ok) // false

	ok, _ = db.GetByID(&product, product1ID)
	log.Println("Get: ", ok) // true ( -> product)

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

	// Reverse (note reversed range)
	c, _ = db.Find(&orders).Reverse().From(mid2012).To(mid2009).Count()
	log.Println("Reverse: ", c) // 3

	// Secondary index on 'customer' - exact index match
	n, _ = db.First(&order).Where("customer", "customer-2").Run()
	log.Println("Index - exact match: ", n) // 1 (-> order )

	// Sum (based on index)
	var sum float64
	db.Find(&orders).Where("shippingfee", 0.00, 10.00).From(mid2009).To(mid2012).Sum(&sum)
	log.Println("Sum: ", sum) // 6.00 (1.00 + 2.00 + 3.00)

	// Secondary index on 'customer' - index range and count
	c, _ = db.Find(&orders).Where("customer", "customer-1", "customer-3").Count()
	log.Println("Index - range: ", c) // 3

	// Secondary index on 'customer' - exact index match, count and date range
	c, _ = db.Find(&orders).Where("customer", "customer-3").From(mid2009).To(time.Now()).Count()
	log.Println("Index - exact match and date range: ", c) // 1

	// Secondary index on 'customer' - index range AND date range
	c, _ = db.Find(&orders).Where("customer", "customer-1", "customer-3").From(mid2009).To(mid2010).Count()
	log.Println("Index - range and date range", c) // 1
```