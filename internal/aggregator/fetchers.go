package aggregator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPUserFetcher struct {
	Client  HTTPClient
	BaseURL string
}

func (fetcher *HTTPUserFetcher) Fetch(ctx context.Context, userID int) (*User, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%d", fetcher.BaseURL, userID), nil)
	res, err := fetcher.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("error fetching user: status code " + fmt.Sprintf("%d", res.StatusCode))
	}

	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

type HTTPPostsFetcher struct {
	Client  HTTPClient
	BaseURL string
}

func (fetcher *HTTPPostsFetcher) Fetch(ctx context.Context) ([]Post, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fetcher.BaseURL, nil)
	res, err := fetcher.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("error fetching posts: status code " + fmt.Sprintf("%d", res.StatusCode))
	}

	var posts []Post
	if err := json.NewDecoder(res.Body).Decode(&posts); err != nil {
		return nil, err
	}
	return posts, nil
}