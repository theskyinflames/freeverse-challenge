package api

import (
	"context"
	"errors"
	"fmt"

	"theskyinflames/freeverse-challenge/internal/app"
	"theskyinflames/freeverse-challenge/internal/domain"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
)

// Product is a DTO
type Product struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Available bool    `json:"available"`
	Price     float64 `json:"price"`
}

var productType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"available": &graphql.Field{
			Type: graphql.Boolean,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
	},
})

// PurchaseResponse is a DTO
type PurchaseResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}

var purchaseResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "PurchaseResponse",
	Fields: graphql.Fields{
		"success": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
		},
		"error": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var purchaseProductInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "PurchaseProductInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"productID": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.ID),
		},
	},
})

func queryType(log cqrs.Logger, bus cqrs.Bus) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"products": &graphql.Field{
				Type:    graphql.NewList(productType),
				Resolve: ProductsResolver(log, bus),
			},
		},
	})
}

func mutationType(log cqrs.Logger, bus cqrs.Bus) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"purchase_product": &graphql.Field{
				Type: purchaseResponseType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(purchaseProductInputType),
					},
				},
				Resolve: PurchaseProductResolver(log, bus),
			},
		},
	})
}

func schema(log cqrs.Logger, bus cqrs.Bus) (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType(log, bus),
		Mutation: mutationType(log, bus),
	})
}

// ExecuteQuery is self-described
func ExecuteQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

// ProductsResolver is a resolver function
func ProductsResolver(log cqrs.Logger, bus cqrs.Bus) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		response, err := bus.Dispatch(context.Background(), app.ProductsQuery{})
		if err != nil {
			log.Printf("something went wrong when executing the query %s: %s\n", app.ProductsQuery{}.Name(), err.Error())
			return nil, nil
		}
		var products []Product
		for _, item := range response.([]app.Product) {
			products = append(products, Product{
				ID:        item.ID.String(),
				Name:      item.Name,
				Available: item.Available,
				Price:     item.Price,
			})
		}
		return products, nil
	}
}

// PurchaseProductResolver is a resolver function
func PurchaseProductResolver(log cqrs.Logger, bus cqrs.Bus) func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
		input, ok := p.Args["input"]
		if !ok {
			log.Printf("input field not found\n")
			return PurchaseResponse{Success: false, Error: errors.New("input field not found").Error()}, nil
		}
		param, ok := input.(map[string]interface{})["productID"]
		if !ok {
			log.Printf("productID field not found\n")
			return PurchaseResponse{Success: false, Error: errors.New("productID field not found").Error()}, nil
		}

		pID, error := uuid.Parse(param.(string))
		if error != nil {
			log.Printf("invalid product UUID\n")
			return PurchaseResponse{Success: false, Error: errors.New("invalid product UUID").Error()}, nil
		}

		_, err := bus.Dispatch(context.Background(), app.PurchaseProductCmd{ID: pID})
		if err != nil {
			log.Printf("something went wrong when executing the command %s: %s\n", app.PurchaseProductCmd{}.Name(), err.Error())
			if errors.Is(err, domain.ErrProductPurchased) {
				return PurchaseResponse{Success: false, Error: errors.New("productID not available for purchasing").Error()}, nil
			}
			if errors.Is(err, app.ErrNotFound) {
				return PurchaseResponse{Success: false, Error: errors.New("productID not found").Error()}, nil
			}
			return PurchaseResponse{Success: false, Error: errors.New("internal error").Error()}, nil
		}

		return PurchaseResponse{Success: true}, nil
	}
}
