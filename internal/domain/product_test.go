package domain_test

import (
	"testing"

	"theskyinflames/graphql-challenge/internal/domain"
	"theskyinflames/graphql-challenge/internal/fixtures"
	"theskyinflames/graphql-challenge/internal/helpers"

	"github.com/stretchr/testify/require"
)

func TestPurchase(t *testing.T) {
	t.Run(`Given a not available product, 
			when it's tried to be purchased, 
			then it returns an error`, func(t *testing.T) {
		p := fixtures.Product{Available: helpers.BoolPtr(false)}.Build()
		require.ErrorIs(t, p.Purchase(), domain.ErrProductPurchased)
	})

	t.Run(`Given an available product, 
			when it's tried to be purchased, 
			then it returns no error`, func(t *testing.T) {
		p := fixtures.Product{Available: helpers.BoolPtr(true)}.Build()
		require.NoError(t, p.Purchase())
		require.Len(t, p.Events(), 1)
	})
}
