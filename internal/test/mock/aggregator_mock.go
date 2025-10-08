package mock

import (
	"bytes"
	"context"
	"go-graphql-aggregator/internal/fetcher"
	"io"
	"net/http"
	"time"
)

var (
	UserMock = &fetcher.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	PostsMock = []fetcher.Post{{UserID: 1}, {UserID: 1}}
)

type MockUserFetcher struct {
	User *fetcher.User
	Err  error
}

// Fetch simulates fetching a user by ID.
func (m *MockUserFetcher) Fetch(ctx context.Context, userID int) (*fetcher.User, error) {
	return m.User, m.Err
}

type MockPostsFetcher struct {
	Posts []fetcher.Post
	Err   error
	Delay time.Duration
}

// Fetch simulates fetching posts by user ID, with an optional delay to test timeouts.
func (m *MockPostsFetcher) Fetch(ctx context.Context, userID int) ([]fetcher.Post, error) {
	if m.Delay > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(m.Delay):
		}
	}
	return m.Posts, m.Err
}

type MockHTTPClient struct {
	response *http.Response
	err      error
}

// Do simulates an HTTP request.
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

// NewMockHTTPClient creates a new MockHTTPClient with the specified response body, status code, and error.
func NewMockHTTPClient(body string, status int, err error) *MockHTTPClient {
	return &MockHTTPClient{
		response: &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
		},
		err: err,
	}
}