package utils

import (
	"net/http"
	"time"
)

// RetryClient wraps an http.Client and adds retry functionality
type RetryClient struct {
	Client     *http.Client
	MaxRetries int
	RetryDelay time.Duration
}

// Do wraps the http.Client's Do method with retry logic
func (c *RetryClient) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= c.MaxRetries; i++ {
		resp, err = c.Client.Do(req)
		if err == nil {
			return resp, nil
		}

		// If this was the last retry, return the error
		if i == c.MaxRetries {
			return nil, err
		}

		// Wait before retrying
		time.Sleep(c.RetryDelay)
	}

	return nil, err
}
