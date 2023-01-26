package api_test

import (
	"context"
	"errors"
	"testing"

	"theskyinflames/graphql-challenge/internal/app"
	"theskyinflames/graphql-challenge/internal/domain"
	"theskyinflames/graphql-challenge/internal/infra/api"

	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/require"
	"github.com/theskyinflames/cqrs-eda/pkg/bus"
)

type busMock struct {
	expectedResult interface{}
	expectedError  error
}

func (bm busMock) Dispatch(context.Context, bus.Dispatchable) (interface{}, error) {
	return bm.expectedResult, bm.expectedError
}

type loggerMock struct {
	calls int
}

func (lm *loggerMock) Printf(string, ...interface{}) {
	lm.calls++
}

func TestProductsResolver(t *testing.T) {
	products := []app.Product{
		{ID: uuid.New(), Name: "product1", Available: true, Price: 1.1},
		{ID: uuid.New(), Name: "product2", Available: true, Price: 2.2},
	}
	testCases := []struct {
		name             string
		bm               busMock
		lm               *loggerMock
		expectedResponse interface{}
		expectedLogCalls int
		expectedError    error
	}{
		{
			name: `Given a bus that returns an error, 
				when it's called, 
				then the error is logged and an empty response is returned`,
			bm: busMock{
				expectedError: errors.New(""),
			},
			lm:               &loggerMock{},
			expectedLogCalls: 1,
		},
		{
			name: `Given a bus that returns a list of products, 
				when it's called, 
				then the the list of products is returned`,
			bm: busMock{
				expectedResult: products,
			},
			lm:               &loggerMock{},
			expectedResponse: products,
		},
	}

	for _, tc := range testCases {
		pr := api.ProductsResolver(tc.lm, tc.bm)
		response, err := pr(graphql.ResolveParams{})
		require.NoError(t, err)
		require.Equal(t, tc.expectedLogCalls, tc.lm.calls)
		if tc.expectedResponse != nil {
			require.Len(t, tc.expectedResponse, len(response.([]api.Product)))
		}
	}
}

func TestPurchaseProductResolver(t *testing.T) {
	randomErr := errors.New("randomErr")
	testCases := []struct {
		name             string
		params           graphql.ResolveParams
		bm               busMock
		lm               *loggerMock
		expectedResponse interface{}
		expectedLogCalls int
		expectedError    error
	}{
		{
			name: `Given a query without input parameter, 
				when it's called, 
				then the error is logged and an empty response is returned`,
			bm: busMock{
				expectedError: randomErr,
			},
			lm:               &loggerMock{},
			expectedLogCalls: 1,
			expectedResponse: api.PurchaseResponse{
				Success: false,
				Error:   "input field not found",
			},
		},
		{
			name: `Given a query without productID, 
				when it's called, 
				then the error is logged and an empty response is returned`,
			params: graphql.ResolveParams{
				Args: map[string]interface{}{
					"input": map[string]interface{}{},
				},
			},
			bm: busMock{
				expectedError: randomErr,
			},
			lm:               &loggerMock{},
			expectedLogCalls: 1,
			expectedResponse: api.PurchaseResponse{
				Success: false,
				Error:   "productID field not found",
			},
		},
		{
			name: `Given a query with and invalid productID, 
				when it's called, 
				then the error is logged and an empty response is returned`,
			params: graphql.ResolveParams{
				Args: map[string]interface{}{
					"input": map[string]interface{}{
						"productID": "invalid",
					},
				},
			},
			bm: busMock{
				expectedError: randomErr,
			},
			lm:               &loggerMock{},
			expectedLogCalls: 1,
			expectedResponse: api.PurchaseResponse{
				Success: false,
				Error:   "invalid product UUID",
			},
		},
		{
			name: `Given a bus that returns an internal error, 
				when it's called, 
				then the error is logged and an empty response is returned`,
			params: graphql.ResolveParams{
				Args: map[string]interface{}{
					"input": map[string]interface{}{
						"productID": uuid.New().String(),
					},
				},
			},
			bm: busMock{
				expectedError: randomErr,
			},
			lm:               &loggerMock{},
			expectedLogCalls: 1,
			expectedResponse: api.PurchaseResponse{
				Success: false,
				Error:   "internal error",
			},
		},
		{
			name: `Given a bus that returns an domain.ErrProductPurchased error, 
				when it's called, 
				then the error is logged and an empty response is returned`,
			params: graphql.ResolveParams{
				Args: map[string]interface{}{
					"input": map[string]interface{}{
						"productID": uuid.New().String(),
					},
				},
			},
			bm: busMock{
				expectedError: domain.ErrProductPurchased,
			},
			lm:               &loggerMock{},
			expectedLogCalls: 1,
			expectedResponse: api.PurchaseResponse{
				Success: false,
				Error:   "productID not available for purchasing",
			},
		},
		{
			name: `Given a bus that returns no error, 
				when it's called, 
				then a success response is returned`,
			params: graphql.ResolveParams{
				Args: map[string]interface{}{
					"input": map[string]interface{}{
						"productID": uuid.New().String(),
					},
				},
			},
			bm: busMock{},
			lm: &loggerMock{},
			expectedResponse: api.PurchaseResponse{
				Success: true,
			},
		},
	}

	for _, tc := range testCases {
		pr := api.PurchaseProductResolver(tc.lm, tc.bm)
		response, err := pr(tc.params)
		require.NoError(t, err, tc.name)
		require.Equal(t, tc.expectedLogCalls, tc.lm.calls, tc.name)
		if tc.expectedResponse != nil {
			require.Equal(t, tc.expectedResponse, response, tc.name)
		}
	}
}
