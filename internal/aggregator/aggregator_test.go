package aggregator_test

import (
	"context"
	"errors"
	"go-graphql-aggregator/internal/aggregator"
	"go-graphql-aggregator/internal/test/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)



func Test_GetUserSummary_Success(t *testing.T){
	assert := assert.New(t)
	userMock := &mock.MockUserFetcher{
		User: &aggregator.User{ID: 1, Name: "John Doe", Email: "john@example.com"},
	}
	postsMock := &mock.MockPostsFetcher{
		Posts: []aggregator.Post{{UserID: 1}, {UserID: 1}, {UserID: 2}},
	}

	agg := aggregator.NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 1)

	assert.Nil(err)
	assert.NotNil(summary)
	assert.Equal("John Doe", summary.Name)
	assert.Equal(2, summary.PostCount)
}

func Test_GetUserSummary_InvalidUserID(t *testing.T){
	assert := assert.New(t)
	userMock := &mock.MockUserFetcher{}
	postsMock := &mock.MockPostsFetcher{}
	agg := aggregator.NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), -1)

	assert.NotNil(err)
	assert.Nil(summary)
	assert.Equal("invalid user ID", err.Error())
}

func Test_GetUserSummary_FetchUserError(t *testing.T){
	assert := assert.New(t)
	userMock := &mock.MockUserFetcher{
		Err: errors.New("user not found"),
	}
	postsMock := &mock.MockPostsFetcher{
		Posts: []aggregator.Post{{UserID: 1}, {UserID: 1}, {UserID: 2}},
	}

	agg := aggregator.NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(summary)
	assert.Equal("failed to fetch user data: user not found", err.Error())
}

func Test_GetUserSummary_FetchPostsError(t *testing.T){
	assert := assert.New(t)
	userMock := &mock.MockUserFetcher{
		User: &aggregator.User{ID: 1, Name: "John Doe", Email: "john@example.com"},
	}
	postsMock := &mock.MockPostsFetcher{
		Err: errors.New("fetch posts error"),
	}

	agg := aggregator.NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(summary)
	assert.Equal("failed to fetch posts data: fetch posts error", err.Error())
}