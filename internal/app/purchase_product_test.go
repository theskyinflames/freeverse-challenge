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

func newInvalidCommand() cqrs.Command {
	return &CommandMock{
		NameFunc: func() string {
			return "invalid_command"
		},
	}
}

func TestPurchaseProduct(t *testing.T) {
	var (
		randomErr = errors.New("")
		product   = fixtures.Product{Available: helpers.BoolPtr(true)}.Build()
	)
	testCases := []struct {
		name              string
		pr                *ProductsRepositoryMock
		cmd               cqrs.Command
		expectedPurchased bool
		expectedErrFunc   func(*testing.T, error)
	}{
		{
			name: `Given an invalid command, when it's called, then an error is returned`,
			cmd:  newInvalidCommand(),
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorAs(t, err, &app.InvalidCommandError{})
			},
		},
		{
			name: `Given a products repository that returns an error on FindByID, 
				when it's called, 
				then an error is returned`,
			cmd: app.PurchaseProductCmd{},
			pr: &ProductsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ uuid.UUID) (domain.Product, error) {
					return domain.Product{}, randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given a products repository that returns an error on Update, 
				when it's called, 
				then an error is returned`,
			cmd: app.PurchaseProductCmd{},
			pr: &ProductsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ uuid.UUID) (domain.Product, error) {
					return product, nil
				},
				UpdateAvailableFunc: func(ctx context.Context, p domain.Product) error {
					return randomErr
				},
			},
			expectedErrFunc: func(t *testing.T, err error) {
				require.ErrorIs(t, err, randomErr)
			},
		},
		{
			name: `Given an available product, 
				when it's purchased, 
				then no error is returned`,
			cmd: app.PurchaseProductCmd{
				ID: uuid.New(),
			},
			pr: &ProductsRepositoryMock{
				FindByIDFunc: func(_ context.Context, _ uuid.UUID) (domain.Product, error) {
					return product, nil
				},
				UpdateAvailableFunc: func(ctx context.Context, p domain.Product) error {
					return nil
				},
			},
		},
	}

	for _, testCase := range testCases {
		ch := app.NewPurchaseProduct(testCase.pr)
		_, err := ch.Handle(context.Background(), testCase.cmd)
		require.Equal(t, testCase.expectedErrFunc == nil, err == nil)
		if err != nil {
			testCase.expectedErrFunc(t, err)
			continue
		}

		require.Len(t, testCase.pr.FindByIDCalls(), 1)
		require.Equal(t, testCase.pr.FindByIDCalls()[0].ID, testCase.cmd.(app.PurchaseProductCmd).ID)
	}
}
