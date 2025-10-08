package aggregator

import (
	"context"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Post struct {
	UserID int `json:"userId"`
}

type UserSummary struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	PostCount int    `json:"postCount"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type UserFetcher interface {
	Fetch(ctx context.Context, userID int) (*User, error)
}

type PostsFetcher interface {
	Fetch(ctx context.Context) ([]Post, error)
}

