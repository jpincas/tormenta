package tormenta

import (
	"github.com/jpincas/gouuidv6"
	"github.com/tinylib/msgp/msgp"
)

func init() {
	// Registering an extension is as simple as matching the
	// appropriate type number with a function that initializes
	// a freshly-allocated object of that type
	msgp.RegisterExtension(99, func() msgp.Extension { return new(gouuidv6.UUID) })
}

//go:generate msgp

// NoModel does not include the Tormenta, so cannot be saved
type NoModel struct {
	SomeData string
}

type Product struct {
	Code          string
	Name          string
	Price         float32
	StartingStock int

	// For a bit more realism and to pad out the serialised data
	Description string
}

type Line struct {
	ProductCode string
	Qty         int
}

type Order struct {
	Model
	Customer int
	Items    []Line
}

const defaultDescription = "On the other hand, we denounce with righteous indignation and dislike men who are so beguiled and demoralized by the charms of pleasure of the moment, so blinded by desire, that they cannot foresee the pain and trouble that are bound to ensue; and equal blame belongs to those who fail in their duty through weakness of will, which is the same as saying through shrinking from toil and pain. These cases are perfectly simple and easy to distinguish. In a free hour, when our power of choice is untrammelled and when nothing prevents our being able to do what we like best, every pleasure is to be welcomed and every pain avoided. But in certain circumstances and owing to the claims of duty or the obligations of business it will frequently occur that pleasures have to be repudiated and annoyances accepted. The wise man therefore always holds in these matters to this principle of selection: he rejects pleasures to secure other greater pleasures, or else he endures pains to avoid worse pains."

// Inventory is a simple product catalogue keyed by product code
var Inventory = map[string]Product{
	"001": Product{"001", "Computer", 999.99, 50, defaultDescription},
	"002": Product{"002", "Mousemat", 9.99, 50, defaultDescription},
	"003": Product{"003", "Mouse", 29.99, 50, defaultDescription},
	"004": Product{"004", "Plant", 6.99, 50, defaultDescription},
	"005": Product{"005", "Desk", 299.99, 50, defaultDescription},
}

// InventoryList is a list version of the catalogue
var InventoryList = getInventoryList()

func getInventoryList() (products []Product) {
	for _, product := range Inventory {
		products = append(products, product)
	}

	return
}
