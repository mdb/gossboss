package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aelsabbahy/goss/outputs"
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

// GetHealthz returns the goss test results served by a goss server
// at its /healthz endpoint.
func (c *Client) GetHealthz(url string) (*outputs.StructuredOutput, error) {
	return c.doRequest(url)
}

func (c *Client) doRequest(url string) (*outputs.StructuredOutput, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	bodyContents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	so := &outputs.StructuredOutput{}
	err = json.Unmarshal(bodyContents, so)
	if err != nil {
		return nil, err
	}

	return so, nil
}
