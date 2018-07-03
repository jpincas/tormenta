# TormentaDB

Simple embedded object persistence for Go - powered by:
- [BadgerDB](https://github.com/dgraph-io/badger)
- [TinyLib MessagePack](https://github.com/tinylib/msgp)
- ['V6' UUIDs](https://github.com/bradleypeabody/gouuidv6)
 
and greatly inspired by [Storm](https://github.com/asdine/storm).

```go
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

	// Get by ID
	nonExistentID := newID()
	product1ID := product1.ID

	var product Product
	ok, _ := db.GetByID(&product, nonExistentID)
	log.Println(ok) // false

	ok, _ = db.GetByID(&product, product1ID)
	log.Println(ok) // true ( -> product)

	// Basic Query
	var products []Product
	n, _ = db.Query(&products).Run()
	log.Println(n) // 2 (-> products)

	// Date range query

	// Make some orders with specific creation times
	var ordersToSave []Tormentable
	dates := []time.Time{
		// Specific years
		time.Date(2009, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2010, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2011, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2012, time.January, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2013, time.January, 1, 1, 0, 0, 0, time.UTC),
	}

	for _, date := range dates {
		ordersToSave = append(ordersToSave, &Order{
			// You wouln't normally do this manually
			// This is just for illustration
			Model: Model{
				ID: gouuidv6.NewFromTime(date),
			},
		})
	}

	// Save the orders
	db.Save(ordersToSave...)

	mid2009 := time.Date(2009, time.June, 1, 1, 0, 0, 0, time.UTC)
	mid2012 := time.Date(2012, time.June, 1, 1, 0, 0, 0, time.UTC)

	var orders []Order
	n, _ = db.Query(&orders).From(mid2009).To(mid2012).Run()
	log.Println(n) // 3 (-> orders )

	// Count only
	// This takes less time than a full query as it:
	// a) skips unmarshalling
	// b) uses Badger's key-only iteration
	count, _ := db.Query(&orders).From(mid2009).To(mid2012).Count()
	log.Println(count) // 3 (-> orders )
}
```