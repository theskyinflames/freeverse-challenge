package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/theskyinflames/cqrs-eda/pkg/cqrs"
)

/* Examples with cURL:
-- get the list of products
curl --request POST \
  --url http://localhost:8080/graphql \
  --header 'Content-Type: application/json' \
  --data '{"query":"{products {id name available price}}"}'


-- purchase a product
curl --request POST \
       --url http://localhost:8080/graphql \
       --header 'Content-Type: application/json' \
       --data '{"query":"mutation {purchase_product(input: {productID: \"ec92361c-3e36-4371-b040-28f608cbe8c6\"}) {success error }}"}'
*/

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

// GraphqlHandler is the HTTP handler for the GraphQL endpoint
func GraphqlHandler(log cqrs.Logger, bus cqrs.Bus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*
			// If it's needed, HTTP headers can be passed to the resolver function in the context
				headers := r.Header
				ctx := context.WithValue(r.Context(), "headers", headers)
		*/

		var p postData
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			w.WriteHeader(400)
			return
		}
		schema, err := schema(log, bus)
		if err != nil {
			log.Printf(fmt.Sprintf("building GraphQL schema: %s\n", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		result := graphql.Do(graphql.Params{
			Context:        r.Context(),
			Schema:         schema,
			RequestString:  p.Query,
			VariableValues: p.Variables,
			OperationName:  p.Operation,
		})

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			fmt.Printf("could not write result to response: %s", err)
		}
	}
}
