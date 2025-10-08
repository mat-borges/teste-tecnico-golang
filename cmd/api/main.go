package main

import (
	"go-graphql-aggregator/internal/aggregator"
	"go-graphql-aggregator/internal/graph"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func newServer() *handler.Server {
	httpClient := http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
		},
	}

	userFetcher := &aggregator.HTTPUserFetcher{
		Client:  &httpClient,
		BaseURL: "https://jsonplaceholder.typicode.com/users",
	}
	postsFetcher := &aggregator.HTTPPostsFetcher{
		Client:  &httpClient,
		BaseURL: "https://jsonplaceholder.typicode.com/posts",
	}

	agg := &aggregator.Aggregator{
		UserFetcher:  userFetcher,
		PostsFetcher: postsFetcher,
		Timeout:      6 * time.Second,
	}

	resolver := &graph.Resolver{Aggregator: agg}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := newServer()

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("Server running at http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
