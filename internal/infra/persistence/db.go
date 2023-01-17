package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

// ConnectToDB opens a connection to the DB
func ConnectToDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open(os.Getenv("DB_DRIVER_NAME"), os.Getenv("DB_URI"))
	if err != nil {
		return nil, fmt.Errorf("something went wrong trying to open database connection: %w", err)
	}
	for { // TODO: adding a back-off retry
		if err := db.PingContext(ctx); err == nil {
			break
		}
	}
	return db, nil
}
