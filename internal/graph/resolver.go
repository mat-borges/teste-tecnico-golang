package graph

import (
	"context"
	"go-graphql-aggregator/internal/aggregator"
	"go-graphql-aggregator/internal/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	Aggregator *aggregator.Aggregator
}

type queryResolver struct{ *Resolver }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver {
    return &queryResolver{r}
}

// UserSummary resolves the userSummary query by fetching and aggregating data.
func (r *queryResolver) UserSummary(ctx context.Context, userID int32) (*model.UserSummary, error) {
	aggSummary, err := r.Aggregator.GetUserSummary(ctx, int(userID))
	if err != nil {
		return nil, err
	}

	modelSummary := &model.UserSummary{
		Name:      aggSummary.Name,
		Email:     aggSummary.Email,
		PostCount: int32(aggSummary.PostCount),
	}
	return modelSummary, nil
}

