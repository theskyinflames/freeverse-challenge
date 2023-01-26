package app_test

import (
	"context"
	"errors"
	"testing"

	"theskyinflames/graphql-challenge/internal/app"
	"theskyinflames/graphql-challenge/internal/domain"
	"theskyinflames/graphql-challenge/internal/fixtures"
	"theskyinflames/graphql-challenge/internal/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
)

func newInvalidQuery() cqrs.Query {
	return &QueryMock{
		NameFunc: func() string {
			return "invalid_query"
		},
	}
}

func TestProducts(t *testing.T) {
	var (
		randomErr      = errors.New("")
		ids            = []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
		names          = []string{"product1", "product2", "product3"}
		prices         = []float64{1.1, 2.2, 3.3}
		availabilities = []bool{true, true, false}
		products       = []domain.Product{
			fixtures.Product{ID: helpers.UUIDPtr(ids[0]), Name: &names[0], Available: &availabilities[0], Price: &prices[0]}.Build(),
			fixtures.Product{ID: helpers.UUIDPtr(ids[1]), Name: &names[1], Available: &availabilities[1], Price: &prices[1]}.Build(),
			fixtures.Product{ID: helpers.UUIDPtr(ids[2]), Name: &names[2], Available: &availabilities[2], Price: &prices[2]}.Build(),
		}
		response = []app.Product{
			{ID: ids[0], Name: names[0], Available: availabilities[0], Price: prices[0]},
			{ID: ids[1], Name: names[1], Available: availabilities[1], Price: prices[1]},
			{ID: ids[2], Name: names[2], Available: availabilities[2], Price: prices[2]},
		}
	)
	testCases := []struct {
		name            string
		pr              *ProductsRepositoryMock
		query           cqrs.Query
		expectedResult  app.ProductsResponse
		expectedErrFunc func(*testing.T, error)
	}{
		{
			name:  `Given an invalid query, when it's called, then an error is returned`,
			query: newInvalidQuery(),
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorAs(t, err, &app.InvalidQueryError{})
			},
		},
		{
			name: `Given a products repository that returns an error on FindAll, 
				when it's called, 
				then an error is returned`,
			query: app.ProductsQuery{},
			pr: &ProductsRepositoryMock{
				FindAllFunc: func(_ context.Context) ([]domain.Product, error) {
					return nil, randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given a products repository that returns a list of products on FindAll, 
				when it's called, 
				then the list of products are returned`,
			query: app.ProductsQuery{},
			pr: &ProductsRepositoryMock{
				FindAllFunc: func(_ context.Context) ([]domain.Product, error) {
					return products, nil
				},
			},
		},
	}

	for _, testCase := range testCases {
		ch := app.NewProducts(testCase.pr)
		result, err := ch.Handle(context.Background(), testCase.query)
		require.Equal(t, testCase.expectedErrFunc == nil, err == nil)
		if err != nil {
			testCase.expectedErrFunc(t, err)
			continue
		}

		require.Len(t, testCase.pr.FindAllCalls(), 1)
		require.Equal(t, response, result)
	}
}
