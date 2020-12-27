package gossboss

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

// NewClient returns a new gossboss Client.
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{},
	}
}

// Healthz represents goss server test results as served by a
// goss server's /healthz endpoint.
type Healthz struct {
	// Result is the /healthz endpoint response body.
	Result *outputs.StructuredOutput

	// URL is the goss server URL.
	URL string

	// Error is any error that was encountered when attempting
	// to fetch the goss test results.
	Error error
}

// CollectHealthzs concurrently retrieves the goss test results from
// each server URL and returns a slice of the results.
func (c *Client) CollectHealthzs(urls []string) []*Healthz {
	ch := make(chan *Healthz)
	results := []*Healthz{}

	for _, url := range urls {
		go c.collectHealthz(url, ch)
	}

	// wait until all goss server test
	// results have been collected.
	for {
		result := <-ch
		results = append(results, result)

		if len(results) == len(urls) {
			close(ch)
			break
		}
	}

	return results
}

// GetHealthz returns the goss test results served by a goss server
// at its /healthz endpoint.
func (c *Client) GetHealthz(url string) (*outputs.StructuredOutput, error) {
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

func (c *Client) collectHealthz(url string, ch chan<- *Healthz) {
	so, err := c.GetHealthz(url)
	ch <- &Healthz{
		Error:  err,
		URL:    url,
		Result: so,
	}
}
