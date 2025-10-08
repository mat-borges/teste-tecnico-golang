package fetcher

import (
	"context"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type UserFetcher interface {
	Fetch(ctx context.Context, userID int) (*User, error)
}

type PostsFetcher interface {
	Fetch(ctx context.Context, userID int) ([]Post, error)
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Post struct {
	UserID int `json:"userId"`
}