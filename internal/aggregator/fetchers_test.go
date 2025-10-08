package aggregator

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func newMockHTTPClient(body string, status int, err error) *mockHTTPClient {
	return &mockHTTPClient{
		response: &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(bytes.NewBufferString(body)),
		},
		err: err,
	}
}

func Test_HTTPUserFetcher_Success(t *testing.T){
	assert := assert.New(t)
	body := `{"id":1,"name":"John Doe","email":"john@example.com"}`
	mockHTTPClient := newMockHTTPClient(body, http.StatusOK, nil)
	mockUserFetcher := &HTTPUserFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/users",
	}

	user, err := mockUserFetcher.Fetch(context.Background(), 1)

	assert.Nil(err)
	assert.NotNil(user)
	assert.Equal(1, user.ID)
	assert.Equal("John Doe", user.Name)
}

func Test_HTTPUserFetcher_RequestError(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := newMockHTTPClient("", http.StatusInternalServerError, errors.New("network error"))
	mockUserFetcher := &HTTPUserFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/users",
	}

	user, err := mockUserFetcher.Fetch(context.Background(), 1)
	assert.NotNil(err)
	assert.Nil(user)
}

func Test_HTTPUserFetcher_InvalidJSON(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := newMockHTTPClient("invalid json", http.StatusOK, nil)
	mockUserFetcher := &HTTPUserFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/users",
	}

	user, err := mockUserFetcher.Fetch(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(user)
}

func Test_HTTPPostsFetcher_Success(t *testing.T){
	assert := assert.New(t)
	body := `[{"userId":1},{"userId":1},{"userId":2}]`
	mockHTTPClient := newMockHTTPClient(body, http.StatusOK, nil)
	mockPostsFetcher := &HTTPPostsFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/posts",
	}

	posts, err := mockPostsFetcher.Fetch(context.Background())

	assert.Nil(err)
	assert.NotNil(posts)
	assert.Len(posts, 3)
	assert.Equal(1, posts[0].UserID)
}

func Test_HTTPPostsFetcher_StatusError(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := newMockHTTPClient("", http.StatusInternalServerError, nil)
	mockPostsFetcher := &HTTPPostsFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/posts",
	}

	posts, err := mockPostsFetcher.Fetch(context.Background())

	assert.NotNil(err)
	assert.Nil(posts)
	assert.Contains(err.Error(), "error fetching posts: status code 500")
}

func Test_HTTPPostsFetcher_RequestError(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := newMockHTTPClient("", http.StatusInternalServerError, errors.New("network error"))
	mockPostsFetcher := &HTTPPostsFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/posts",
	}

	posts, err := mockPostsFetcher.Fetch(context.Background())

	assert.NotNil(err)
	assert.Nil(posts)
}

func Test_HTTPPostsFetcher_InvalidJSON(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := newMockHTTPClient("invalid json", http.StatusOK, nil)
	mockPostsFetcher := &HTTPPostsFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/posts",
	}

	posts, err := mockPostsFetcher.Fetch(context.Background())

	assert.NotNil(err)
	assert.Nil(posts)
}