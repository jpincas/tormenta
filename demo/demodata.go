package demo

import (
	"math/rand"
	"time"

	"github.com/jpincas/gouuidv6"

	tormenta "github.com/jpincas/tormenta/tormentadb"
)

// PopulateWithDemoData fills the provided DB with data
func PopulateWithDemoData(db *tormenta.DB, n int) {
	// Generate demo data
	orders := Orders(n)
	products := Products(n)

	// Save it
	db.Save(orders...)
	db.Save(products...)
}

// Orders creates n demo orders
func Orders(n int) (orders []tormenta.Tormentable) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < n; i++ {
		order := Order{
			Model: tormenta.Model{
				// Random creation date from now to a year ago
				ID: gouuidv6.NewFromTime(
					randomDate(time.Date(2016, time.June, 1, 23, 0, 0, 0, time.UTC)),
				),
			},
			Customer:                "a customer",
			Department:              1,
			ShippingFee:             4.99,
			ContainsProhibitedItems: false,
		}

		orders = append(orders, &order)
	}

	return
}

// Products creates n demo products
func Products(n int) (products []tormenta.Tormentable) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < n; i++ {
		product := Product{
			Code:          "SKU000",
			Name:          "Product Name",
			Price:         9.99,
			StartingStock: 50,
			Tags:          []string{"tag1", "tag2"},
			Departments:   []int{1, 2, 3},
			Description:   DefaultDescription,
		}

		products = append(products, &product)
	}

	return
}

// randomDate generates a random date between the given 'from' time and now
func randomDate(from time.Time) time.Time {
	f := from.Unix()
	diff := time.Now().Unix() - f
	r := rand.Int63n(diff)
	seconds := f + r
	return time.Unix(seconds, 0)
}
