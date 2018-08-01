package demo

import (
	"errors"

	"github.com/jpincas/gouuidv6"
	tormenta "github.com/jpincas/tormenta/tormentadb"
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
	tormenta.Model
	Code          string
	Name          string `tormenta:"split"`
	Price         float64
	StartingStock int
	Tags          []string
	Departments   []int

	// For a bit more realism and to pad out the serialised data
	Description string `tormenta:"noindex"`
}

type Line struct {
	ProductCode string
	Qty         int
}

type Order struct {
	tormenta.Model
	Customer                string
	Department              int
	ShippingFee             float64
	HasShipped              bool
	Items                   []Line
	ContainsProhibitedItems bool `msg:"-" tormenta:"noindex"`
	OrderSaved              bool `msg:"-" tormenta:"noindex"`
	OrderRetrieved          bool `msg:"-" tormenta:"noindex"`
}

func (o Order) PreSave() error {
	if o.ContainsProhibitedItems {
		return errors.New("Cannot place order - contains prohibited items")
	}

	return nil
}

func (o *Order) PostSave() {
	o.OrderSaved = true
}

func (o *Order) PostGet() {
	o.OrderRetrieved = true
}

const DefaultDescription = "On the other hand, we denounce with righteous indignation and dislike men who are so beguiled and demoralized by the charms of pleasure of the moment, so blinded by desire, that they cannot foresee the pain and trouble that are bound to ensue; and equal blame belongs to those who fail in their duty through weakness of will, which is the same as saying through shrinking from toil and pain. These cases are perfectly simple and easy to distinguish. In a free hour, when our power of choice is untrammelled and when nothing prevents our being able to do what we like best, every pleasure is to be welcomed and every pain avoided. But in certain circumstances and owing to the claims of duty or the obligations of business it will frequently occur that pleasures have to be repudiated and annoyances accepted. The wise man therefore always holds in these matters to this principle of selection: he rejects pleasures to secure other greater pleasures, or else he endures pains to avoid worse pains."

// Inventory is a simple product catalogue keyed by product code
var Inventory = map[string]Product{
	"001": Product{
		Code:          "001",
		Name:          "Computer",
		Price:         999.99,
		StartingStock: 50,
		Description:   DefaultDescription},
	"002": Product{
		Code:          "002",
		Name:          "Mouse",
		Price:         9.99,
		StartingStock: 50,
		Description:   DefaultDescription},
	"003": Product{
		Code:          "003",
		Name:          "Mousemat",
		Price:         5.99,
		StartingStock: 50,
		Description:   DefaultDescription},
	"004": Product{
		Code:          "004",
		Name:          "Desk",
		Price:         199.99,
		StartingStock: 50,
		Description:   DefaultDescription},
	"005": Product{
		Code:          "005",
		Name:          "Plant",
		Price:         4.99,
		StartingStock: 50,
		Description:   DefaultDescription},
}

// InventoryList is a list version of the catalogue
var InventoryList = getInventoryList()

func getInventoryList() (products []Product) {
	for _, product := range Inventory {
		products = append(products, product)
	}

	return
}
