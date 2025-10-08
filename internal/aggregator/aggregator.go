package aggregator

import (
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