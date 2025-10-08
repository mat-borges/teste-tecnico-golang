package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ---------------- USERS -------------------

type HTTPUserFetcher struct {
	Client  HTTPClient
	BaseURL string
}

// Fetch fetches user data by userID.
func (fetcher *HTTPUserFetcher) Fetch(ctx context.Context, userID int) (*User, error) {
	req, err  := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%d", fetcher.BaseURL, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("creating user request: %w", err)
	}

	var lastErr error
	for attempt := range 3 {
		res, err := fetcher.Client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("doing user request: %w", err)
		} else {
			defer res.Body.Close()
			if res.StatusCode == http.StatusOK {
				var user User
				if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
					return nil, fmt.Errorf("decoding user response: %w", err)
				}
				return &user, nil
			}
				lastErr = fmt.Errorf("fetching user: status code %d", res.StatusCode)
		}
		wait := time.Duration(1<<attempt) * 100 * time.Millisecond
		select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(wait):
        }

	}
	return nil, lastErr
}

// ---------------- POSTS -------------------

type HTTPPostsFetcher struct {
	Client  HTTPClient
	BaseURL string
}

// Fetch fetches posts by userID.
func (fetcher *HTTPPostsFetcher) Fetch(ctx context.Context, userID int) ([]Post, error) {
	u, err := url.Parse(fetcher.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid posts base url: %w", err)
	}

	var lastErr error
	for attempt := range 3 {
		q := u.Query()
		if userID > 0 {
			q.Set("userId", fmt.Sprintf("%d", userID))
		}
		u.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			lastErr = fmt.Errorf("creating posts request: %w", err)
		} else {
			res, err := fetcher.Client.Do(req)
			if err != nil {
				lastErr = fmt.Errorf("doing posts request: %w", err)
			} else {
				defer res.Body.Close()

				if res.StatusCode == http.StatusOK {
					var posts []Post
					if err := json.NewDecoder(res.Body).Decode(&posts); err != nil {
						return nil, fmt.Errorf("decoding posts response: %w", err)
					}
					return posts, nil
				}
				lastErr = fmt.Errorf("fetching posts: status code %d", res.StatusCode)
			}
		}
		wait := time.Duration(1<<attempt) * 100 * time.Millisecond
		select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(wait):
        }
	}
	return nil, lastErr
}