package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
)

// Product is a DTO
type Product struct {
	ID        uuid.UUID
	Name      string
	Available bool
	Price     float64
}

// ProductsResponse is a DTO
type ProductsResponse []Product

// ProductsQuery is a query
type ProductsQuery struct{}

// ProductsName is self-described
var ProductsName = "products"

// Name implements Query interface
func (q ProductsQuery) Name() string {
	return ProductsName
}

// Products is a query handler
type Products struct {
	pr ProductsRepository
}

// NewProducts is a constructor
func NewProducts(pr ProductsRepository) Products {
	return Products{pr: pr}
}

// Handle implements the QueryHandler interface
func (qh Products) Handle(ctx context.Context, query cqrs.Query) (cqrs.QueryResult, error) {
	_, ok := query.(ProductsQuery)
	if !ok {
		return nil, NewInvalidQueryError(ProductsName, query.Name())
	}

	p, err := qh.pr.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []Product
	for _, item := range p {
		response = append(response, Product{
			ID:        item.ID(),
			Name:      item.Name(),
			Available: item.Available(),
			Price:     item.Price(),
		})
	}

	return response, nil
}
