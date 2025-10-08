package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewServer_InitializesCorrectly(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv := newServer(ctx)
	assert.NotNil(srv, "server should be created successfully")

	req := httptest.NewRequest("GET", "/query", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	assert.NotEqual(500, w.Code)
}

func Test_NewServer_ContextCancelled(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	srv := newServer(ctx)
	assert.Nil(srv, "server should be nil if context is cancelled")
}

func Test_Main_StartupAndShutdown(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	srv := newServer(ctx)
	assert.NotNil(srv, "server should initialize")

	handler := http.NewServeMux()
	handler.Handle("/query", srv)

	testServer := httptest.NewServer(handler)
	defer testServer.Close()

	// requisição real
	resp, err := http.Get(testServer.URL + "/query")
	assert.Nil(err, "request should not return error")
	assert.NotNil(resp)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	assert.NotEmpty(body, "response body should not be empty")

	testServer.CloseClientConnections()
	time.Sleep(100 * time.Millisecond)
}