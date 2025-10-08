package main

import (
	"go-graphql-aggregator/internal/aggregator"
	"go-graphql-aggregator/internal/graph"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	userFetcher := &aggregator.HTTPUserFetcher{
		Client:  http.DefaultClient,
		BaseURL: "https://jsonplaceholder.typicode.com/users",
	}

	postsFetcher := &aggregator.HTTPPostsFetcher{
		Client:  http.DefaultClient,
		BaseURL: "https://jsonplaceholder.typicode.com/posts",
	}

	aggregatorInstance := &aggregator.Aggregator{
		UserFetcher:  userFetcher,
		PostsFetcher: postsFetcher,
	}
	resolver := &graph.Resolver{
		Aggregator: aggregatorInstance,
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))


	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("Server running at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
