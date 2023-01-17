package app

import (
	"context"

	"theskyinflames/freeverse-challenge/internal/domain"

	"github.com/google/uuid"
)

//go:generate moq -stub -out zmock_app_repositories_test.go -pkg app_test . ProductsRepository

// ProductsRepository is self-described
type ProductsRepository interface {
	FindByID(ctx context.Context, ID uuid.UUID) (domain.Product, error)
	FindAll(ctx context.Context) ([]domain.Product, error)
	UpdateAvailable(ctx context.Context, p domain.Product) error
}
