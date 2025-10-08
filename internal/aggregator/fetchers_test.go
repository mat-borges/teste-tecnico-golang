package aggregator_test

import (
	"context"
	"errors"
	"go-graphql-aggregator/internal/aggregator"
	"go-graphql-aggregator/internal/test/mock"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)


func Test_HTTPUserFetcher_Success(t *testing.T){
	assert := assert.New(t)
	body := `{"id":1,"name":"John Doe","email":"john@example.com"}`
	mockHTTPClient := mock.NewMockHTTPClient(body, http.StatusOK, nil)
	mockUserFetcher := &aggregator.HTTPUserFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/users",
	}

	user, err := mockUserFetcher.Fetch(context.Background(), 1)

	assert.Nil(err)
	assert.NotNil(user)
	assert.Equal(1, user.ID)
	assert.Equal("John Doe", user.Name)
}

func Test_HTTPUserFetcher_StatusError(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := mock.NewMockHTTPClient("", http.StatusInternalServerError, nil)
	mockUserFetcher := &aggregator.HTTPUserFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/users",
	}

	user, err := mockUserFetcher.Fetch(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(user)
	assert.Contains(err.Error(), "error fetching user: status code 500")
}

func Test_HTTPUserFetcher_RequestError(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := mock.NewMockHTTPClient("", http.StatusInternalServerError, errors.New("network error"))
	mockUserFetcher := &aggregator.HTTPUserFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/users",
	}

	user, err := mockUserFetcher.Fetch(context.Background(), 1)
	assert.NotNil(err)
	assert.Nil(user)
}

func Test_HTTPUserFetcher_InvalidJSON(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := mock.NewMockHTTPClient("invalid json", http.StatusOK, nil)
	mockUserFetcher := &aggregator.HTTPUserFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/users",
	}

	user, err := mockUserFetcher.Fetch(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(user)
}

func Test_HTTPUserFetcher_CreateRequestError(t *testing.T) {
	assert := assert.New(t)
	mockUserFetcher := &aggregator.HTTPUserFetcher{
		Client:  http.DefaultClient,
		BaseURL: "://invalid-url",
	}

	user, err := mockUserFetcher.Fetch(context.Background(), 1)
	assert.NotNil(err)
	assert.Nil(user)
	assert.Contains(err.Error(), "error creating user request")
}

// ---------------- POSTS -------------------

func Test_HTTPPostsFetcher_Success(t *testing.T){
	assert := assert.New(t)
	body := `[{"userId":1},{"userId":1}]`
	mockHTTPClient := mock.NewMockHTTPClient(body, http.StatusOK, nil)
	mockPostsFetcher := &aggregator.HTTPPostsFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/posts",
	}

	posts, err := mockPostsFetcher.Fetch(context.Background(), 1)

	assert.Nil(err)
	assert.NotNil(posts)
	assert.Len(posts, 2)
	assert.Equal(1, posts[0].UserID)
}

func Test_HTTPPostsFetcher_StatusError(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := mock.NewMockHTTPClient("", http.StatusInternalServerError, nil)
	mockPostsFetcher := &aggregator.HTTPPostsFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/posts",
	}

	posts, err := mockPostsFetcher.Fetch(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(posts)
	assert.Contains(err.Error(), "error fetching posts: status code 500")
}

func Test_HTTPPostsFetcher_RequestError(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := mock.NewMockHTTPClient("", http.StatusInternalServerError, errors.New("network error"))
	mockPostsFetcher := &aggregator.HTTPPostsFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/posts",
	}

	posts, err := mockPostsFetcher.Fetch(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(posts)
}

func Test_HTTPPostsFetcher_InvalidJSON(t *testing.T){
	assert := assert.New(t)
	mockHTTPClient := mock.NewMockHTTPClient("invalid json", http.StatusOK, nil)
	mockPostsFetcher := &aggregator.HTTPPostsFetcher{
		Client:  mockHTTPClient,
		BaseURL: "http://example.com/posts",
	}

	posts, err := mockPostsFetcher.Fetch(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(posts)
}

func Test_HTTPPostsFetcher_InvalidBaseURL(t *testing.T) {
	assert := assert.New(t)
	fetcher := &aggregator.HTTPPostsFetcher{
		Client:  http.DefaultClient,
		BaseURL: "://bad-url",
	}

	posts, err := fetcher.Fetch(context.Background(), 1)
	println(err)
	assert.NotNil(err)
	assert.Nil(posts)
	assert.Contains(err.Error(), "invalid posts base url")
}