package mock

import (
	"bytes"
	"context"
	"go-graphql-aggregator/internal/aggregator"
	"io"
	"net/http"
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
}

func (m *MockPostsFetcher) Fetch(ctx context.Context) ([]aggregator.Post, error) {
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