package aggregator

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Aggregator struct {
	httpClient 	HTTPClient
	apiBaseURL 	string
	timeout     time.Duration
}

func New(httpClient HTTPClient, apiBaseURL string, timeout time.Duration) *Aggregator {
	return &Aggregator{
		httpClient: httpClient,
		apiBaseURL: apiBaseURL,
		timeout:    timeout,
	}
}

func (agg *Aggregator) GetUserSummary(ctx context.Context, userID int) (*UserSummary, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	ctx, cancel := context.WithTimeout(ctx, agg.timeout)
	defer cancel()

	var(
		wg 	 	sync.WaitGroup
		user 	*User
		posts 	[]Post
		errUser, errPosts error
	)

	wg.Add(2)

	go func(){
		defer wg.Done()
		u, err := FetchUser(ctx, agg.httpClient, agg, agg.apiBaseURL+"/users", userID)
		if err != nil {
			errUser = err
			return
		}
		user = u
	}()

	go func(){
		defer wg.Done()
		p, err := FetchPosts(ctx, agg.httpClient, agg.apiBaseURL+"/posts")
		if err != nil {
			errPosts = err
			return
		}
		posts = p
	}()

	wg.Wait()

	if errUser != nil {
		return nil, errors.New("failed to fetch user data: " + errUser.Error())
	}
	if errPosts != nil {
		return nil, errors.New("failed to fetch posts data: " + errPosts.Error())
	}

	count := 0
	for _, post := range posts {
		if post.UserID == userID {
			count++
		}
	}

	return &UserSummary{Name: user.Name, Email: user.Email, PostCount: count}, nil
}