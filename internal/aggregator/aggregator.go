package aggregator

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

type Aggregator struct {
	UserFetcher  UserFetcher
	PostsFetcher PostsFetcher
	Timeout      time.Duration
}

// NewAggregator creates a new Aggregator instance.
func NewAggregator(userFetcher UserFetcher, postsFetcher PostsFetcher, timeout time.Duration) *Aggregator {
	return &Aggregator{
		UserFetcher:  userFetcher,
		PostsFetcher: postsFetcher,
		Timeout:      timeout,
	}
}

// GetUserSummary fetches user data and their posts, then aggregates the information.
func (agg *Aggregator) GetUserSummary(ctx context.Context, userID int) (*UserSummary, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	ctx, cancel := context.WithTimeout(ctx, agg.Timeout)
	defer cancel()

	var user *User
	var posts []Post

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		u, err := agg.UserFetcher.Fetch(ctx, userID)
		if err != nil {
			return fmt.Errorf("fetching user: %w", err)
		}
		user = u
		return nil
	})

	g.Go(func() error {
		p, err := agg.PostsFetcher.Fetch(ctx, userID)
		if err != nil {
			return fmt.Errorf("fetching posts: %w", err)
		}
		posts = p
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &UserSummary{
		Name:      user.Name,
		Email:     user.Email,
		PostCount: len(posts),
	}, nil
}