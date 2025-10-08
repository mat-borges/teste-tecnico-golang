package graph_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-graphql-aggregator/internal/aggregator"
	"go-graphql-aggregator/internal/graph"
	"go-graphql-aggregator/internal/test/mock"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/stretchr/testify/assert"
)


func Test_UserSummaryQuery_Success(t *testing.T) {
	assert := assert.New(t)
	mockAgg := &aggregator.Aggregator{
		UserFetcher: &mock.MockUserFetcher{
			User: &aggregator.User{ID: 1, Name: "John Doe", Email: "john@example.com"},
		},
		PostsFetcher: &mock.MockPostsFetcher{
			Posts: []aggregator.Post{{UserID: 1}, {UserID: 1}, {UserID: 2}},
		},
		Timeout: 0,
	}

	resolver := &graph.Resolver{Aggregator: mockAgg}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.POST{})

	body := `{"query": "query { userSummary(userId: 1) { name email postCount } }"}`
	req := httptest.NewRequest("POST", "/query", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	var resp struct {
		Data struct {
			UserSummary struct {
				Name      string
				Email     string
				PostCount int32
			}
		}
		Errors []any
	}
	err := json.NewDecoder(w.Body).Decode(&resp)

	assert.Nil(err)

	assert.Equal(200, w.Code)
	assert.NotNil(resp.Data.UserSummary)
	assert.Equal("John Doe", resp.Data.UserSummary.Name)
	assert.Equal("john@example.com", resp.Data.UserSummary.Email)
	assert.Equal(int32(2), resp.Data.UserSummary.PostCount)
}

func Test_UserSummaryQuery_UserFetcherError(t *testing.T) {
	assert := assert.New(t)

	mockAgg := &aggregator.Aggregator{
		UserFetcher: &mock.MockUserFetcher{
			Err: errors.New("user fetch failed"),
		},
		PostsFetcher: &mock.MockPostsFetcher{
			Posts: []aggregator.Post{{UserID: 1}, {UserID: 1}, {UserID: 2}},
		},
	}

	resolver := &graph.Resolver{Aggregator: mockAgg}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.POST{})

	body := `{"query": "query { userSummary(userId: 1) { name email postCount } }"}`
	req := httptest.NewRequest("POST", "/query", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	assert.Equal(200, w.Code)

	var resp struct {
		Data   map[string]any
		Errors []struct {
			Message string
		}
	}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.Nil(err)
	assert.Len(resp.Errors, 1)
	assert.Contains(resp.Errors[0].Message, "failed to fetch user data: user fetch failed")
	assert.Nil(resp.Data["userSummary"])
}

func Test_UserSummaryQuery_PostsFetcherError(t *testing.T) {
	assert := assert.New(t)

	mockAgg := &aggregator.Aggregator{
		UserFetcher: &mock.MockUserFetcher{
			User: &aggregator.User{ID: 1, Name: "John Doe", Email: "john@example.com"},
		},
		PostsFetcher: &mock.MockPostsFetcher{Err: errors.New("posts fetch failed")},
	}

	resolver := &graph.Resolver{Aggregator: mockAgg}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.POST{})

	body := `{"query": "query { userSummary(userId: 1) { name email postCount } }"}`
	req := httptest.NewRequest("POST", "/query", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)
	t.Log("Response:", w.Body.String())

	assert.Equal(200, w.Code)

	var resp struct {
		Errors []struct {
			Message string
		}
	}
	_ = json.NewDecoder(w.Body).Decode(&resp)

	assert.Len(resp.Errors, 1)
	assert.Contains(resp.Errors[0].Message, "failed to fetch posts data")
}

func Test_UserSummaryQuery_InvalidUserID(t *testing.T) {
	assert := assert.New(t)

	mockAgg := &aggregator.Aggregator{
		UserFetcher: &mock.MockUserFetcher{},
		PostsFetcher: &mock.MockPostsFetcher{},
	}

	resolver := &graph.Resolver{Aggregator: mockAgg}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	srv.AddTransport(transport.POST{})

	body := `{"query": "query { userSummary(userId: -1) { name email postCount } }"}`
	req := httptest.NewRequest("POST", "/query", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	srv.ServeHTTP(w, req)

	assert.Equal(200, w.Code)

	var resp struct {
		Data   map[string]any
		Errors []struct {
			Message string
		}
	}
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.Nil(err)
	assert.Len(resp.Errors, 1)
	assert.Contains(resp.Errors[0].Message, "invalid user ID")
	assert.Nil(resp.Data["userSummary"])
}
