package main

import (
	"context"
	"go-graphql-aggregator/internal/aggregator"
	"go-graphql-aggregator/internal/config"
	"go-graphql-aggregator/internal/fetcher"
	"go-graphql-aggregator/internal/graph"
	"go-graphql-aggregator/internal/logger"
	"go-graphql-aggregator/internal/middleware"
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

func newServer(ctx context.Context, cfg *config.Config) *handler.Server {
	httpClient := http.Client{
		Timeout: cfg.HTTPTimeout,
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

	userFetcher := &fetcher.HTTPUserFetcher{
		Client:  &httpClient,
		BaseURL: cfg.UsersBaseURL,
	}
	postsFetcher := &fetcher.HTTPPostsFetcher{
		Client:  &httpClient,
		BaseURL: cfg.PostsBaseURL,
	}

	agg := &aggregator.Aggregator{
		UserFetcher:  userFetcher,
		PostsFetcher: postsFetcher,
		Timeout:      cfg.AggTimeout,
	}

	select {
	case <-ctx.Done():
		logger.Log.Error("Server initialization cancelled", "error", ctx.Err())
		return nil
	default:
	}

	resolver := &graph.Resolver{Aggregator: agg}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	if os.Getenv("ENABLE_INTROSPECTION") == "1" {
		srv.Use(extension.Introspection{})
	}
	if os.Getenv("ENABLE_APQ") == "1" {
		srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})
	}

	return srv
}

func main() {
	logger.Init()
	logger.Log.Info("server starting...")

	startupCtx, startupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startupCancel()

	cfg := config.LoadConfig()

	srvHandler := newServer(startupCtx, cfg)
	if srvHandler == nil {
		logger.Log.Error("server initialization failed")
		os.Exit(1)
	}

	http.Handle("/", middleware.LoggingAndRecoveryMiddleware(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", middleware.LoggingAndRecoveryMiddleware(srvHandler))

	httpServer := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      nil,
	}

	serverErrCh := make(chan error, 1)

	go func() {
		logger.Log.Info("starting server", "port", cfg.ServerPort)
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
		logger.Log.Info("received OS signal, initiating shutdown", "signal", sig)
	case err := <-serverErrCh:
		if err != nil {
			logger.Log.Error("server runtime error", "error", err)
		} else {
			logger.Log.Info("server stopped normally")
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("error during server shutdown", "error", err)
	} else {
		logger.Log.Info("server shutdown completed")
	}
}
