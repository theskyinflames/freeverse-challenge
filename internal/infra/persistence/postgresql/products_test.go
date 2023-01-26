//go:build test_db
// +build test_db

package postgresql_test

import (
	"context"
	"database/sql"
	"testing"
	"theskyinflames/graphql-challenge/internal/fixtures"
	"theskyinflames/graphql-challenge/internal/helpers"
	"theskyinflames/graphql-challenge/internal/infra/persistence"
	"theskyinflames/graphql-challenge/internal/infra/persistence/postgresql"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	ctx := context.Background()

	// Starting a PostgreSQL container for the tests
	// This container will be destroyed when the test suite finishes..
	dbConfig := NewConfig()
	CreateContainerizedPgDB(t, dbConfig)

	// Connect to the DB
	db, err := sql.Open("postgres", dbConfig.DatabaseURL())
	require.NoError(t, err)
	defer func() {
		_ = db.Close()
	}()

	// Running DB migrations
	require.NoError(t, persistence.RunMigrations(ctx, db, "postgres", "file://./migrations"))

	// Running the tests suite
	suite.Run(t, &PostgreSQLTestSuite{db: db})
}

// Define the tests suite
type PostgreSQLTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (suite *PostgreSQLTestSuite) TestFindByID() {
	t := suite.T()

	// Insert fixture data into DB
	var (
		id        = uuid.New()
		name      = "product1"
		available = true
		price     = 1.1
	)
	_, err := suite.db.Exec(
		"INSERT INTO products (id, name, price, available) VALUES ($1, $2, $3, $4)",
		id,
		name,
		price,
		available,
	)
	require.NoError(t, err)

	pr := postgresql.NewProductsRepository(suite.db)
	found, err := pr.FindByID(context.Background(), id)
	require.NoError(t, err)

	require.Equal(t, id, found.ID())
	require.Equal(t, available, found.IsAvailable())
}

func (suite *PostgreSQLTestSuite) TestFindAll() {
	t := suite.T()
	pr := postgresql.NewProductsRepository(suite.db)
	found, err := pr.FindAll(context.Background())
	require.NoError(t, err)

	require.True(t, len(found) > 0)
}

func (suite *PostgreSQLTestSuite) TestUpdateAvailable() {
	t := suite.T()

	// Insert fixture data into DB
	var (
		id        = uuid.New()
		name      = "product22"
		available = false
		price     = 1.1
	)
	_, err := suite.db.Exec(
		"INSERT INTO products (id, name, price, available) VALUES ($1, $2, $3, $4)",
		id,
		name,
		price,
		available,
	)
	require.NoError(t, err)

	p := fixtures.Product{ID: &id, Name: &name, Available: helpers.BoolPtr(true), Price: &price}.Build()
	pr := postgresql.NewProductsRepository(suite.db)
	require.NoError(t, pr.UpdateAvailable(context.Background(), p))

	found, err := pr.FindByID(context.Background(), id)
	require.NoError(t, err)

	require.True(t, found.IsAvailable())
}
