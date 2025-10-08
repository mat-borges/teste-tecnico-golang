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
	userMock := &mock.MockUserFetcher{User: mock.UserMock}
	postsMock := &mock.MockPostsFetcher{Posts: mock.PostsMock}

	agg := aggregator.NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 1)

	assert.Nil(err)
	assert.NotNil(summary)
	assert.Equal("John Doe", summary.Name)
	assert.Equal("john@example.com", summary.Email)
	assert.Equal(2, summary.PostCount)
}

func Test_GetUserSummary_InvalidUserID(t *testing.T){
	assert := assert.New(t)
	agg := aggregator.NewAggregator(&mock.MockUserFetcher{}, &mock.MockPostsFetcher{}, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 0)

	assert.NotNil(err)
	assert.Nil(summary)
	assert.Contains(err.Error(), "invalid user ID")
}

func Test_GetUserSummary_FetchUserError(t *testing.T){
	assert := assert.New(t)
	userMock := &mock.MockUserFetcher{Err: errors.New("user not found")}
	postsMock := &mock.MockPostsFetcher{}

	agg := aggregator.NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(summary)
	assert.Equal("fetching user: user not found", err.Error())
}

func Test_GetUserSummary_FetchPostsError(t *testing.T){
	assert := assert.New(t)
	userMock := &mock.MockUserFetcher{User: mock.UserMock}
	postsMock := &mock.MockPostsFetcher{Err: errors.New("fetch posts error")}
	agg := aggregator.NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(summary)
	assert.Equal("fetching posts: fetch posts error", err.Error())
}

func TestAggregator_GetUserSummary_Timeout(t *testing.T) {
	assert := assert.New(t)
	userMock := &mock.MockUserFetcher{User: mock.UserMock}
	postsMock := &mock.MockPostsFetcher{Delay: 3 * time.Second}
	agg := aggregator.NewAggregator(userMock, postsMock, 1*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(summary)
	assert.Contains(err.Error(), "context deadline exceeded")
}