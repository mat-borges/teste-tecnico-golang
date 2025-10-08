package aggregator

import (
	"context"
	"fmt"
	"go-graphql-aggregator/internal/fetcher"
	"go-graphql-aggregator/internal/logger"
	"time"

	"golang.org/x/sync/errgroup"
)

type Aggregator struct {
	UserFetcher  fetcher.UserFetcher
	PostsFetcher fetcher.PostsFetcher
	Timeout      time.Duration
}

// NewAggregator creates a new Aggregator instance.
func NewAggregator(userFetcher fetcher.UserFetcher, postsFetcher fetcher.PostsFetcher, timeout time.Duration) *Aggregator {
	return &Aggregator{
		UserFetcher:  userFetcher,
		PostsFetcher: postsFetcher,
		Timeout:      timeout,
	}
}

// GetUserSummary fetches user details and their posts, returning a summary.
func (agg *Aggregator) GetUserSummary(ctx context.Context, userID int) (*UserSummary, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}

	start := time.Now()

	if agg.Timeout <= 0 {
		agg.Timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, agg.Timeout)
	defer cancel()

	var user *fetcher.User
	var posts []fetcher.Post

	g, ctx := errgroup.WithContext(ctx)

	logger.Log.Info("fetch start", "userId", userID, "timeout", agg.Timeout)

	g.Go(func() error {
		u, err := agg.UserFetcher.Fetch(ctx, userID)
		if err != nil {
			logger.Log.Error("fetch user failed", "userId", userID, "error", err)
			return fmt.Errorf("fetching user: %w", err)
		}
		user = u
		logger.Log.Info("fetch user done", "userId", userID)
		return nil
	})

	g.Go(func() error {
		p, err := agg.PostsFetcher.Fetch(ctx, userID)
		if err != nil {
			logger.Log.Error("fetch posts failed", "userId", userID, "error", err)
			return fmt.Errorf("fetching posts: %w", err)
		}
		posts = p
		logger.Log.Info("fetch posts done", "userId", userID, "count", len(p))
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Log.Error("aggregation failed", "userId", userID, "error", err)
		return nil, err
	}

	elapsed := time.Since(start)
	logger.Log.Info("aggregation complete",
		"userId", userID,
		"posts", len(posts),
		"elapsed_ms", elapsed.Milliseconds(),
	)

	return &UserSummary{
		Name:      user.Name,
		Email:     user.Email,
		PostCount: len(posts),
	}, nil
}