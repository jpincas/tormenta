package demo

import (
	"fmt"
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
			Department:              rand.Intn(20),
			ShippingFee:             float64(rand.Intn(20)) + .99,
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
			Code:          fmt.Sprintf("SKU00%v", i),
			Name:          fmt.Sprintf("Product %v", i),
			Price:         float64(rand.Intn(100)) + 9.99,
			StartingStock: rand.Intn(100),
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
