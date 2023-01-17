# Freeverse.io challenge

This a coding challenge to implement the freeverse-challenge service requisites:
```
  Instructions:
    Write a Go program that exposes an endpoint (using either RPC, REST, or GraphQL) that allows a user to purchase a product. The products are stored in a SQL database and they are unique.

    The endpoint should accept the following parameters:
    * productID: The ID of the product to be purchased

    The endpoint should return the following response:
    * success: A boolean value indicating whether the purchase was successful
    * error: A string describing the error that occurred, if any

    The endpoint should perform the following actions:
    1. Retrieve the product with the specified productID from the SQL database
    2. Check if the product is available for purchase
    3. If the product is not available, return an error
    4. If the product is available, set the available field to false and return a success response

  Requirements:
    * The Go program must be written using idiomatic Go and should follow the Go coding standards
    * The SQL schema should be created as you see fit and should be accessed using a Go library (e.g. database/sql)
    * The Go program should be well-documented and easy to understand
    * The endpoint should be thoroughly tested using a testing framework (e.g. testing)
    * The project should be dockerized
    * Try to minimize the use of third-party libraries
```

## Architectural decisions taken

To write this service I've taken into account these aspects:

* Of course [Go](https://go.dev) is the programming language I've used.

* The IT team told me that they use GraphQL APIs, so, this is the API I've decided to implement. You will find the [GraphQL schema here](./schema/products.graphql)

```graphql
    type Product {
      id: String!
      name: String!
      price: Float!
      available: Boolean!
    }

    type Query {
      products: [Product!]!
    }

    type PurchaseResponse {
      success: Boolean!
      error: String
    }

    input PurchaseProductInput {
      productID: String!
    }

    type Mutation {
      purchaseProduct(input: PurchaseProductInput!): PurchaseResponse
    }
```

* Given that the Architecture guidelines, if I understood properly during my visit to your office, are still not defined, I've tried to apply the most common architectural patterns for Service oriented environments:
  * [Hexagonal Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))
  * [CQRS](https://learn.microsoft.com/es-es/azure/architecture/patterns/cqrs) to implement the application services
  * Command-Bus to dispatch the CQRS commands from the graphql resolvers
  * [Domain events](https://dev.to/isaacojeda/ddd-cqrs-aplicando-domain-events-en-aspnet-core-o6n)
  * Events-Bus

* I've added unit tests to all packages.
* I've applied [SOLID principles](https://en.wikipedia.org/wiki/SOLID) also
* I use [databse migrations](https://www.prisma.io/dataguide/types/relational/what-are-database-migrations) for initialize the data model by creating the products table and filling some fixture data. These migrations are executed at PostgreSQL docker starting time.j

## How to run it

This is dockerized service. Tu start it only do:

```sh
  make docker-run
  make docker-logs
```

If you take a look to the [docker-compose.yml](./docker-compose.yml) you'll see that the service starts two containers:

* PosgreSQL container as persistence layer provider

* Service container  

## How to try it

These are the GraphQL requests that the service's API provides:

  * A Query to get the list of products. Needed to know the IDs of the products to be purchase.

  ```sh
    curl --request POST \
      --url http://localhost:8080/graphql \
      --header 'Content-Type: application/json' \
      --data '{"query":"{products {id name available price}}"}'
  ```

  * A Mutation to purchase products.
  
  ```sh
    curl --request POST \
       --url http://localhost:8080/graphql \
       --header 'Content-Type: application/json' \
       --data '{"query":"mutation {purchase_product(input: {productID: \"ec92361c-3e36-4371-b040-28f608cbe8c6\"}) {success error }}"}'
  ```

## How to test it

The service includes unit tests. They can be run this way:

```sh
  make test-unit
  make test-db
```

## Repo layout

* schema - GraphQL schema implemented
* scripts - an script to compile the Docker container locally, for development purposes
* cmd - where the *main.go* lives
* internal - used to [reduce the public API surface](https://dave.cheney.net/2019/10/06/use-internal-packages-to-reduce-your-public-api-surface)
* internal/app - CQRS layer, application services
* internal/domain - where the domain entities and business rules lives
* internal/fixtures builders - needed fixtures for the tests
* internal/helpers - misc helpers used to improved the code reading
* internal/infra - infrastructure layer
* internal/infra/persistence - storage service. Implements *repository* pattern
* internal/infra/persistence/postgres - Database migrations and repository implementation
* internal/infra/api - GraphQL API

There are also other files used for development purposes:

  * internal/golangci.yml - used as linter
  * internal/revive.toml - used as linter
  * Makefile

## Tooling and libs used

To implement the solution I've used:

* Go libs:
  * Third party:
      * github.com/go-chi/chi v1.5.4
	    * github.com/golang-migrate/migrate/v4 v4.15.2
	    * github.com/google/uuid v1.3.0
	    * github.com/graphql-go/graphql v0.8.0
	    * github.com/lib/pq v1.10.0
	    * github.com/ory/dockertest/v3 v3.9.1
	    * github.com/rs/cors v1.8.3
	    * github.com/stretchr/testify v1.8.1
  * My own libs:
	    * github.com/theskyinflames/cqrs-eda v1.2.5
* Tooling:
  * MacOS Ventura 13.1
  * Go 1.19.4
  * Docker 20.10.21
  * PostgreSQL Docker Image latest
