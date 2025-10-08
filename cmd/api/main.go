package main

import (
	"context"
	"go-graphql-aggregator/internal/aggregator"
	"go-graphql-aggregator/internal/graph"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func newServer(ctx context.Context) *handler.Server {
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

	select {
	case <-ctx.Done():
		log.Printf("Server initialization cancelled: %v", ctx.Err())
		return nil
	default:
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
	startupCtx, startupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startupCancel()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srvHandler := newServer(startupCtx)
	if srvHandler == nil {
		log.Fatal("server initialization failed (cancelled or error)")
	}

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srvHandler)

	httpServer := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      nil,
	}

	serverErrCh := make(chan error, 1)

	go func() {
		log.Printf("Starting server at http://localhost:%s/ ...", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		} else {
			serverErrCh <- nil
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("Received OS signal: %v — initiating graceful shutdown", sig)
	case err := <-serverErrCh:
		if err != nil {
			log.Printf("Server runtime error: %v", err)
		} else {
			log.Printf("Server stopped normally")
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	} else {
		log.Printf("Server shutdown completed")
	}
}
