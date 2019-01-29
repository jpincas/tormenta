## Update Jan 2019

I've refactored the repo and remove dthe REST layer and GUI for now.  Going forward, I'll put those in separate repos so that this one is just for the core Tormenta persistence engine, which is currently my priority.  The REST layer will make a reappearance some time in the future, along with a better GUI built on top of it (in Elm!).

Another big change is the move away from MessagePack in favour of good old JSON.  Although this will worsen performance, it means no more codegen, so easier to get started, plus it means I can do some pretty cool 'pass straight through to JSON' stuff which might be useful if you're building a JSON api and don't need to do any intermediate processing.

Still WIP, but being used in two projects, both in dev and going to production in the next couple of months, so the API will have to become at least somewhat stable soon.

Currently working on improving the query engine to at least be able to run multiple 'where' clause queries.

Any questions, hit me up.

# TormentaDB (WIP)

Tormenta is a functionality layer over BadgerDB key/value store.  It provides simple, embedded object persistence for Go projects with some limited data querying capabilities and ORM-like features.  It uses date-tted IDs so is particuarly good for data sets that are natually chronological, like financial transactions, soical media posts etc. Powered by:

- [BadgerDB](https://github.com/dgraph-io/badger)
- ['V6' UUIDs](https://github.com/bradleypeabody/gouuidv6)
 
and greatly inspired by [Storm](https://github.com/asdine/storm).

## Features

- Uses good old JSON for persistence (using `jsoniter` for better performance that the stdlib)
- Date-tted unique IDs - get date range querying and created at field 'for free'
- Simple basic API for saving and retrieving your objects
- Automatic secondary indexing on all fields (can be skipped)
- Option to index by individual words in strings
- More complex, ORM-like querying including exact matches, text prefix, ranges, reverse, limit and offset
- Fast counts and sums using Badger's 'key only' iteration
- Business logic using 'triggers' on save and get


## Quick How To

- Add import `"github.com/jpincas/tormenta"`
- Add `tormenta.Model` to structs you want to persist 
- Add `tormenta:"noindex"` tag to fields you want to exclude from secondary indexing
- Add `tormenta:"split"` tag to string fields where you'd like to index each word separately instead of the the whole sentence
- Open a DB connection with `db, err := tormenta.Open("mydatadirectory")` (dont forget to `defer db.Close()`). For auto-deleting test DB, use `tormenta.OpenTest`
- Save a single entity with `db.Save(&MyEntity)` or multiple entities with `db.Save(&MyEntity1, &MyEntity2)`.
- Get a single entity by ID with `db.Get(&MyEntity, entityID)`.
- Construct a query to find single or mutliple entities with `db.First(&MyEntity)` or `db.Find(&MyEntities)` respectively. 
- Build up the query by chaining methods: `From()/.To()` to add a date range, `Match("indexName", value)` to add an exact match index search, `Range("indexname", start, end)` to add a range search, `StartsWith("indexname", "prefix")` for a text prefix search, `.Reverse()` to reverse the fullStruct of searching/results and `.Limit()/.Offset()` to limit the number of results.
- Kick off the query with `.Run()`, or `.Count()` if you just need the count.  `.QuickSum()` is also available for float/int index searches.
- Add business logic by specifying `.PreSave()`, `.PostSave()` and `.PostGet()` methods on your structs.
	
## Example

```go
package example

import (
	"fmt"
	"log"
	"time"

	"https://github.com/bradleypeabody/gouuidv6"
	"github.com/jpincas/tormenta"
)

// Define your data.
// Include tormenta.Model to get date tted IDs, last updated field etc
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
	db, _ := tormenta.OpenTest("data/tests", tormenta.DefaultOptions)
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
	log.Println("Basic - date range: ", n) // 3 (-> fullStructs )

	// First
	n, _ = db.First(&fullStruct).From(mid2009).To(mid2012).Run()
	log.Println("First - found: ", n) // 1 (-> fullStruct )

	// First (not found)
	n, _ = db.First(&fullStruct).From(time.Now()).To(time.Now()).Run()
	log.Println("First - not found: ", n) // 0

	// Count only (fast!)
	c, _ := db.Find(&fullStructs).From(mid2009).To(mid2012).Count()
	log.Println("Count: ", c) // 3

	// Limit
	n, _ = db.Find(&fullStructs).From(mid2009).To(mid2012).Limit(2).Run()
	log.Println("Limit: ", n) // 2

	// Offset
	n, _ = db.Find(&fullStructs).From(mid2009).To(mid2012).Limit(2).Offset(1).Run()
	log.Println("Limit and offset: ", n) // 2 (same count, different results)

	// Reverse
	c, _ = db.Find(&fullStructs).Reverse().From(mid2009).To(mid2012).Count()
	log.Println("Reverse: ", c) // 3

	// Secondary index on 'customer' - exact index match
	n, _ = db.First(&fullStruct).Match("customer", "customer-2").Run()
	log.Println("Index - exact match: ", n) // 1 (-> fullStruct )

	// Secondary index on 'customer' - prefix match
	n, _ = db.First(&fullStruct).StartsWith("customer", "customer-").Run()
	log.Println("Index - prefix match: ", n) // 5 (-> fullStruct )

	// QuickSum (based on index)
	var sum float64
	db.Find(&fullStructs).Range("shippingfee", 0.00, 10.00).From(mid2009).To(mid2012).QuickSum(&sum)
	log.Println("QuickSum: ", sum) // 6.00 (1.00 + 2.00 + 3.00)

	// Secondary index on 'customer' - index range and count
	c, _ = db.Find(&fullStructs).Range("customer", "customer-1", "customer-3").Count()
	log.Println("Index - range: ", c) // 3

	// Secondary index on 'customer' - exact index match, count and date range
	c, _ = db.Find(&fullStructs).Match("customer", "customer-3").From(mid2009).To(time.Now()).Count()
	log.Println("Index - exact match and date range: ", c) // 1

	// Secondary index on 'customer' - index range AND date range
	c, _ = db.Find(&fullStructs).Range("customer", "customer-1", "customer-3").From(mid2009).To(mid2010).Count()
	log.Println("Index - range and date range", c) // 1
}
```

## To Do

- [x] Delete
- [x] Logic triggers (preSave, postSave, postGet)
- [ ] Easy joins
- [ ] Arbitrary filter functions for indexes
- [ ] Better error reporting from query construction
- [ ] Better protection against unsupported types being passed around as interfaces
- [ ] Fully benchmarked simulation of a real-world use case
- [x] Slices as indexes -> multiple index entries
- [ ] Rebuild indices command
- [ ] Stack multiple queries, execute as AND/OR, execute in parallel
- [ ] For large iterations, how could we run parts in parallel? Would have to specify with a tag?
- [x] Split-string indexing with 'split' tag
- [x] 'Starts with' index match
- [x] Indexes on by default
- [ ] Multiple entity `Get()`
- [ ] Bulk unmarshall rather than 1 at a time? Concurrent?
- [ ] JSON dump/ backup