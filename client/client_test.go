package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockServer(path, body string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if path != r.RequestURI {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "not found")
			return
		}

		w.WriteHeader(200)
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

func TestGetResults(t *testing.T) {
	endpoint := "/healthz"
	respStr := gossResponse()
	server := mockServer(endpoint, respStr)
	defer server.Close()

	c := NewClient()

	resp, err := c.GetResults(server.URL + endpoint)
	if err != nil {
		t.Error("GetResults should not error")
	}

	if resp.StatusCode != 200 {
		t.Errorf("GetResults should return a '200' status code; got '%v'", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("GetResults should return a readable response body")
	}

	r := string(bodyBytes)
	if r != respStr {
		t.Errorf("GetResults should return '%s'; got '%s'", respStr, r)
	}
}
