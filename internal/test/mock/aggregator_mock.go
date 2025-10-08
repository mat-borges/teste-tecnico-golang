package mock

import (
	"bytes"
	"context"
	"go-graphql-aggregator/internal/aggregator"
	"io"
	"net/http"
	"time"
)

var (
	UserMock = &aggregator.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	PostsMock = []aggregator.Post{{UserID: 1}, {UserID: 1}}
)

type MockUserFetcher struct {
	User *aggregator.User
	Err  error
}

func (m *MockUserFetcher) Fetch(ctx context.Context, userID int) (*aggregator.User, error) {
	return m.User, m.Err
}

type MockPostsFetcher struct {
	Posts []aggregator.Post
	Err   error
	Delay time.Duration
}

func (m *MockPostsFetcher) Fetch(ctx context.Context, userID int) ([]aggregator.Post, error) {
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

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func NewMockHTTPClient(body string, status int, err error) *MockHTTPClient {
	return &MockHTTPClient{
		response: &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
		},
		err: err,
	}
}