package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchUser(ctx context.Context, client HTTPClient, agg *Aggregator, baseURL string, userID int) (*User, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%d", baseURL, userID), nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func FetchPosts(ctx context.Context, client HTTPClient, postsURL string) ([]Post, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, postsURL, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error fetching posts: status code %d", res.StatusCode)
	}

	var posts []Post
	if err := json.NewDecoder(res.Body).Decode(&posts); err != nil {
		return nil, err
	}
	return posts, nil
}