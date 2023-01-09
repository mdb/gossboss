package gossboss

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aelsabbahy/goss/outputs"
)

// Client is an HTTP client for interacting with goss servers.
type Client struct {
	// HTTPClient is a *http.Client.
	HTTPClient *http.Client
}

// NewClient returns a new gossboss goss server Client.
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{},
	}
}

// Healthzs represents an aggregate collection of goss server test results
// across multiple goss servers' /healthz endpoints.
type Healthzs struct {
	Healthzs []*Healthz
	Summary  *Summary
}

// Summary is a summary of all Healthzs results.
type Summary struct {
	Failed  int `json:"failed-count"`
	Errored int `json:"errored-count"`
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
func (c *Client) CollectHealthzs(urls []string) *Healthzs {
	hzs := &Healthzs{
		Summary: &Summary{
			Failed:  0,
			Errored: 0,
		},
	}
	ch := make(chan *Healthz, len(urls))

	for _, url := range urls {
		go c.collectHealthz(url, ch)
	}

	// wait until all goss server test
	// results have been collected.
	for len(hzs.Healthzs) < len(urls) {
		hz := <-ch
		hzs.Healthzs = append(hzs.Healthzs, hz)

		if hz.Error == nil && hz.Result.Summary.Failed != 0 {
			hzs.Summary.Failed += hz.Result.Summary.Failed
		}

		if hz.Error != nil {
			hzs.Summary.Errored++
		}
	}

	return hzs
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
