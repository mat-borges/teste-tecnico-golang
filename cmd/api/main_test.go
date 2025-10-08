package main

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewServer_InitializesCorrectly(t *testing.T) {
	assert := assert.New(t)
	srv := newServer()
	assert.NotNil(srv)

	req := httptest.NewRequest("GET", "/query", nil)
	w := httptest.NewRecorder()
	srv.ServeHTTP(w, req)

	assert.NotEqual(500, w.Code, "Server should not return 500 on GET request")
}

