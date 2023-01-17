package persistence

import (
	"context"
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // migrations run as part of service starting
)

// RunMigrations updates the data schema in the persistence layer
func RunMigrations(ctx context.Context, db *sql.DB, dbName, path string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(path, dbName, driver)
	if err != nil {
		return err
	}
	return m.Up()
}
