package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockServer(path, body string, responseCode int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if path != r.RequestURI {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "not found")
			return
		}

		w.WriteHeader(responseCode)
		fmt.Fprintf(w, body)
	}))

	return server
}

func gossResponse() string {
	return `{
  	"results": [{
      "duration": 48248740,
      "err": null,
      "expected": [
        "true"
      ],
      "found": [
        "true"
      ],
      "human": "",
      "meta": null,
      "property": "reachable",
      "resource-id": "tcp://some-server.com:443",
      "resource-type": "Addr",
      "result": 0,
      "successful": true,
      "summary-line": "Addr: tcp://some-server.com:443: reachable: matches expectation: [true]",
      "test-type": 0,
      "title": ""
    }],
		"summary": {
			"failed-count": 0,
			"summary-line": "Count: 1, Failed: 0, Duration: 0.048s",
			"test-count": 1,
			"total-duration": 48441102
		}
	}`
}

func TestGetHealthz(t *testing.T) {
	endpoint := "/healthz"
	respStr := gossResponse()

	server := mockServer(endpoint, respStr, 200)
	defer server.Close()

	c := NewClient()

	resp, err := c.GetHealthz(server.URL + endpoint)
	if err != nil {
		t.Error("GetHealthz should not error")
	}

	if resp.Summary.Failed != 0 {
		t.Errorf("GetHealthz should return StructuredOuput reporting a Summary.Failed of '0'; got '%v'", resp.Summary.Failed)
	}
}

func TestCollectHealthz(t *testing.T) {
	type response struct {
		code int
		body string
	}

	tests := []struct {
		name      string
		responses []response
	}{{
		name: "all 3 servers respond 200",
		responses: []response{{
			code: 200,
			body: gossResponse(),
		}, {
			code: 200,
			body: gossResponse(),
		}, {
			code: 200,
			body: gossResponse(),
		},
		}}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			endpoint := "/healthz"

			servers := []*httptest.Server{}
			for _, r := range test.responses {
				s := mockServer(endpoint, r.body, r.code)
				defer s.Close()
				servers = append(servers, s)
			}

			c := NewClient()

			serverURLs := []string{}
			for _, s := range servers {
				serverURLs = append(serverURLs, s.URL+endpoint)
			}

			resps := c.CollectAllHealthz(serverURLs)

			if len(resps) != len(serverURLs) {
				t.Errorf("CollectAllHealthz should return results from '%v' servers; got '%v'", len(serverURLs), len(resps))
			}

			if resps[0].URL != serverURLs[0] && resps[0].URL != serverURLs[1] && resps[0].URL != serverURLs[2] {
				t.Error("CollectAllHealthz should return a slice of Healthz, each reporting a URL")
			}

			if resps[0].Result.Summary.Failed != 0 {
				t.Errorf("CollectAllHealthz should return a slice of Healthz, each reporting a Result.Summary.Failed of '0'; got '%v'", resps[0].Result.Summary.Failed)
			}

			if resps[0].Error != nil {
				t.Errorf("CollectAllHealthz should return a slice of Healthz, each reporting a nil Error; got '%v'", resps[0].Error)
			}
		})
	}
}
