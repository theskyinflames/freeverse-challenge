package domain

import (
	"errors"

	"github.com/google/uuid"
	"github.com/theskyinflames/cqrs-eda/pkg/ddd"
)

// Product is an entity
type Product struct {
	ddd.AggregateBasic

	name      string
	available bool
	// float64 type used for simplicity, but in a real service a lib like better use github.com/shopspring/decimal
	// take a look to this article to understand which is the problem with using float64 for currency values:
	//    https://waclawthedev.medium.com/inaccurate-float32-and-float64-how-to-avoid-the-trap-in-go-golang-6de59e66aed9
	price float64
}

// NewProduct is a constructor
func NewProduct(ID uuid.UUID, name string, price float64) Product {
	return Product{
		AggregateBasic: ddd.NewAggregateBasic(ID),
		name:           name,
		price:          price,
	}
}

// Name is a getter
func (p Product) Name() string {
	return p.name
}

// Available is a getter
func (p Product) Available() bool {
	return p.available
}

// Price is a getter
func (p Product) Price() float64 {
	return p.price
}

// ErrProductPurchased is self-described
var ErrProductPurchased = errors.New("product not available for purchasing")

// Purchase tries to purchase the product
func (p *Product) Purchase() error {
	if !p.available {
		return ErrProductPurchased
	}
	p.available = false

	p.RecordEvent(NewProductPurchasedEvent(*p))
	return nil
}

// IsPurchased is self-described
func (p Product) IsPurchased() bool {
	return !p.available
}

// IsAvailable is self-described
func (p Product) IsAvailable() bool {
	return p.available
}

// Hydrate hydrates a product instance. It's used to retrieve entities from DB.
func (p *Product) Hydrate(ID uuid.UUID, name string, available bool, price float64) {
	p.AggregateBasic = ddd.NewAggregateBasic(ID)
	p.name = name
	p.available = available
	p.price = price
}
