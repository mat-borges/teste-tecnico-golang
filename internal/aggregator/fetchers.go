package aggregator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPUserFetcher struct {
	client  HTTPClient
	baseURL string
}

func (fetcher *HTTPUserFetcher) Fetch(ctx context.Context, userID int) (*User, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%d", fetcher.baseURL, userID), nil)
	res, err := fetcher.client.Do(req)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

type HTTPPostsFetcher struct {
	client  HTTPClient
	baseURL string
}

func (fetcher *HTTPPostsFetcher) Fetch(ctx context.Context) ([]Post, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fetcher.baseURL, nil)
	res, err := fetcher.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return nil, errors.New("error fetching posts: status code " + fmt.Sprintf("%d", res.StatusCode))
	}

	var posts []Post
	if err := json.NewDecoder(res.Body).Decode(&posts); err != nil {
		return nil, err
	}
	return posts, nil
}