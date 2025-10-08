package aggregator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockUserFetcher struct {
	User *User
	Err  error
}

func (m *MockUserFetcher) Fetch(ctx context.Context, userID int) (*User, error) {
	return m.User, m.Err
}

type MockPostsFetcher struct {
	Posts []Post
	Err   error
}

func (m *MockPostsFetcher) Fetch(ctx context.Context) ([]Post, error) {
	return m.Posts, m.Err
}

func Test_GetUserSummary_Success(t *testing.T){
	assert := assert.New(t)
	userMock := &MockUserFetcher{
		User: &User{ID: 1, Name: "John Doe", Email: "john@example.com"},
	}
	postsMock := &MockPostsFetcher{
		Posts: []Post{{UserID: 1}, {UserID: 1}, {UserID: 2}},
	}

	agg := NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 1)

	assert.Nil(err)
	assert.NotNil(summary)
	assert.Equal("John Doe", summary.Name)
	assert.Equal(2, summary.PostCount)
}

func Test_GetUserSummary_InvalidUserID(t *testing.T){
	assert := assert.New(t)
	userMock := &MockUserFetcher{}
	postsMock := &MockPostsFetcher{}
	agg := NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), -1)

	assert.NotNil(err)
	assert.Nil(summary)
	assert.Equal("invalid user ID", err.Error())
}

func Test_GetUserSummary_FetchUserError(t *testing.T){
	assert := assert.New(t)
	userMock := &MockUserFetcher{
		Err: errors.New("user not found"),
	}
	postsMock := &MockPostsFetcher{
		Posts: []Post{{UserID: 1}, {UserID: 1}, {UserID: 2}},
	}

	agg := NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(summary)
	assert.Equal("failed to fetch user data: user not found", err.Error())
}

func Test_GetUserSummary_FetchPostsError(t *testing.T){
	assert := assert.New(t)
	userMock := &MockUserFetcher{
		User: &User{ID: 1, Name: "John Doe", Email: "john@example.com"},
	}
	postsMock := &MockPostsFetcher{
		Err: errors.New("fetch posts error"),
	}

	agg := NewAggregator(userMock, postsMock, 2*time.Second)
	summary, err := agg.GetUserSummary(context.Background(), 1)

	assert.NotNil(err)
	assert.Nil(summary)
	assert.Equal("failed to fetch posts data: fetch posts error", err.Error())
}