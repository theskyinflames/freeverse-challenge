package main

import (
	"context"
	"fmt"
	"os"

	"theskyinflames/graphql-challenge/cmd/service"
	"theskyinflames/graphql-challenge/internal/infra/persistence"
	"theskyinflames/graphql-challenge/internal/infra/persistence/postgresql"
)

const srvPort = ":80"

func main() {
	ctx := context.Background()
	db, err := persistence.ConnectToDB(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer func() {
		_ = db.Close()
	}()

	fmt.Printf("db migration run starting\n")
	if err := persistence.RunMigrations(ctx, db, os.Getenv("DB_NAME"), os.Getenv("DB_MIGRATIONS_PATH")); err != nil {
		fmt.Printf("something went wrong trying to run database migrations: %s\n", err.Error())
		os.Exit(-1)
	}
	fmt.Printf("db migration run finished\n")

	service.Run(context.Background(), srvPort, postgresql.NewProductsRepository(db))
}
