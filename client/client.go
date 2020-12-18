package client

import (
	"net/http"
)

// Client is an HTTP client for interacting with goss servers.
type Client struct {
	HTTPClient *http.Client
}

// NewClient returns a new Client with the URL it's passed.
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{},
	}
}

// GetResults returns the goss test results served by a goss server.
func (c *Client) GetResults(url string) (*http.Response, error) {
	return c.doRequest(url)
}

func (c *Client) doRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}
