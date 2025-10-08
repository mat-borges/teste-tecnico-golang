package aggregator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type HTTPUserFetcher struct {
	Client  HTTPClient
	BaseURL string
}

func (fetcher *HTTPUserFetcher) Fetch(ctx context.Context, userID int) (*User, error) {
	req, err  := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%d", fetcher.BaseURL, userID), nil)
	if err != nil {
		return nil, errors.New("error creating user request: " + err.Error())
	}
	res, err := fetcher.Client.Do(req)
	if err != nil {
		return nil, errors.New("error doing user request: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("error fetching user: status code " + fmt.Sprintf("%d", res.StatusCode))
	}

	var user User
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, errors.New("error decoding user response: " + err.Error())
	}
	return &user, nil
}

type HTTPPostsFetcher struct {
	Client  HTTPClient
	BaseURL string
}

func (fetcher *HTTPPostsFetcher) Fetch(ctx context.Context, userID int) ([]Post, error) {
	u, err := url.Parse(fetcher.BaseURL)
	if err != nil {
		return nil, errors.New("invalid posts base url: " + err.Error())
	}

	q := u.Query()
	if userID > 0 {
		q.Set("userId", fmt.Sprintf("%d", userID))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.New("error creating posts request: " + err.Error())
	}

	res, err := fetcher.Client.Do(req)
	if err != nil {
		return nil, errors.New("error doing posts request: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("error fetching posts: status code " + fmt.Sprintf("%d", res.StatusCode))
	}

	var posts []Post
	if err := json.NewDecoder(res.Body).Decode(&posts); err != nil {
		return nil, errors.New("error decoding posts response: " + err.Error())
	}
	return posts, nil
}