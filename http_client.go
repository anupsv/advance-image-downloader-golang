package main

import (
	"net/http"
	"time"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type StandardHTTPClient struct {
	client *http.Client
}

func NewStandardHTTPClient() *StandardHTTPClient {
	return &StandardHTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *StandardHTTPClient) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}
