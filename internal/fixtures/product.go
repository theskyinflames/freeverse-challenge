package fixtures

import (
	"theskyinflames/freeverse-challenge/internal/domain"

	"github.com/google/uuid"
)

// Product is fixture
type Product struct {
	ID        *uuid.UUID
	Name      *string
	Available *bool
	Price     *float64
}

// Build is self-described
func (e Product) Build() domain.Product {
	id := uuid.New()
	if e.ID != nil {
		id = *e.ID
	}
	name := "product1"
	if e.Name != nil {
		name = *e.Name
	}
	var available bool
	if e.Available != nil {
		available = *e.Available
	}
	price := 1.1
	if e.Price != nil {
		price = *e.Price
	}

	p := domain.Product{}
	p.Hydrate(id, name, available, price)
	return p
}
