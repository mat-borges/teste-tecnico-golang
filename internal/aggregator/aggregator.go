package aggregator

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Aggregator struct {
	UserFetcher  UserFetcher
	PostsFetcher PostsFetcher
	Timeout      time.Duration
}

func NewAggregator(userFetcher UserFetcher, postsFetcher PostsFetcher, timeout time.Duration) *Aggregator {
	return &Aggregator{
		UserFetcher:  userFetcher,
		PostsFetcher: postsFetcher,
		Timeout:      timeout,
	}
}

func (agg *Aggregator) GetUserSummary(ctx context.Context, userID int) (*UserSummary, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	ctx, cancel := context.WithTimeout(ctx, agg.Timeout)
	defer cancel()

	var(
		wg 	 	sync.WaitGroup
		user 	*User
		posts 	[]Post
		errOnce sync.Once
		errVar error
	)

	setErr := func(err error) {
		errOnce.Do(func() {
			errVar = errors.New(err.Error())
			cancel()
		})
	}

	wg.Add(2)

	go func(){
		defer wg.Done()
		u, err := agg.UserFetcher.Fetch(ctx, userID)
		if err != nil {
			setErr(fmt.Errorf("failed to fetch user! %w", err))
			return
		}
		user = u
	}()

	go func(){
		defer wg.Done()
		p, err := agg.PostsFetcher.Fetch(ctx, userID)
		if err != nil {
			setErr(fmt.Errorf("failed to fetch posts! %w", err))
			return
		}
		posts = p
	}()

	wg.Wait()

	if errVar != nil {
		return nil, errVar
	}

	postCount := len(posts)

	return &UserSummary{
		Name:      user.Name,
		Email:     user.Email,
		PostCount: postCount,
	}, nil
}