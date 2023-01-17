package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"theskyinflames/freeverse-challenge/internal/app"
	"theskyinflames/freeverse-challenge/internal/infra/api"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

// Run Starts the API server
func Run(ctx context.Context, srvPort string, pr app.ProductsRepository) {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
	})
	r.Use(cors.Handler)
	r.Use(middleware.Logger)

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log := log.New(os.Stdout, "freeverse-challenge: ", os.O_APPEND)

	bus := app.BuildCommandQueryBus(log, app.BuildEventsBus(), pr)
	r.Post("/graphql", api.GraphqlHandler(log, bus))

	fmt.Printf("serving at port %s\n", srvPort)
	if err := http.ListenAndServe(srvPort, r); err != nil {
		fmt.Printf("something went wrong trying to start the server: %s\n", err.Error())
	}
}
