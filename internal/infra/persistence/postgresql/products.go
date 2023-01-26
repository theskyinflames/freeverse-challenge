package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"theskyinflames/graphql-challenge/internal/app"
	"theskyinflames/graphql-challenge/internal/domain"

	"github.com/google/uuid"
)

// ProductsRepository is a repository
type ProductsRepository struct {
	db *sql.DB
}

// NewProductsRepository is a constructor
func NewProductsRepository(db *sql.DB) ProductsRepository {
	return ProductsRepository{db: db}
}

// FindByID is a finder
func (pr ProductsRepository) FindByID(ctx context.Context, ID uuid.UUID) (domain.Product, error) {
	stmt, err := pr.db.Prepare("SELECT id,name,available,price FROM products WHERE id=$1")
	if err != nil {
		return domain.Product{}, err
	}

	var (
		foundID   uuid.UUID
		name      string
		available bool
		optPrice  sql.NullFloat64
	)

	err = stmt.QueryRow(ID).Scan(&foundID, &name, &available, &optPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Product{}, app.ErrNotFound
		}
		return domain.Product{}, err
	}

	var p domain.Product
	p.Hydrate(foundID, name, available, price(optPrice))
	return p, nil
}

// FindAll is a finder
func (pr ProductsRepository) FindAll(ctx context.Context) ([]domain.Product, error) {
	rows, err := pr.db.QueryContext(ctx, "SELECT id,name,available,price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var (
			foundID   uuid.UUID
			name      string
			available bool
			optPrice  sql.NullFloat64
		)

		err = rows.Scan(&foundID, &name, &available, &optPrice)
		if err != nil {
			return nil, err
		}

		var p domain.Product
		p.Hydrate(foundID, name, available, price(optPrice))
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// UpdateAvailable updates the available field of the product
func (pr ProductsRepository) UpdateAvailable(ctx context.Context, p domain.Product) error {
	result, err := pr.db.Exec("UPDATE products set available=$1 WHERE ID=$2", p.IsAvailable(), p.ID())
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update product: rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("update product: %w", errors.New("not found"))
	}

	return nil
}

func price(opt sql.NullFloat64) float64 {
	price := 0.0
	if opt.Valid {
		price = opt.Float64
	}
	return price
}
