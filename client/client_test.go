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

func gossResponse(isSuccess bool) string {
	failedCount := 0
	if !isSuccess {
		failedCount = 1
	}

	return fmt.Sprintf(`{
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
      "successful": %t,
      "summary-line": "Addr: tcp://some-server.com:443: reachable: matches expectation: [true]",
      "test-type": 0,
      "title": ""
    }],
		"summary": {
			"failed-count": %v,
			"summary-line": "Count: 1, Failed: 0, Duration: 0.048s",
			"test-count": 1,
			"total-duration": 48441102
		}
	}`, isSuccess, failedCount)
}

func TestGetHealthz(t *testing.T) {
	endpoint := "/healthz"
	respStr := gossResponse(true)

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
		name        string
		failedCount int
		errorCount  int
		responses   []response
	}{{
		name:        "all servers respond 200",
		failedCount: 0,
		errorCount:  0,
		responses: []response{{
			code: 200,
			body: gossResponse(true),
		}, {
			code: 200,
			body: gossResponse(true),
		}}}, {
		name:        "1 server responds 500",
		failedCount: 1,
		errorCount:  0,
		responses: []response{{
			code: 500,
			body: gossResponse(false),
		}, {
			code: 200,
			body: gossResponse(true),
		}},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			endpoint := "/healthz"

			servers := []*httptest.Server{}
			serverURLs := []string{}
			for _, r := range test.responses {
				s := mockServer(endpoint, r.body, r.code)
				defer s.Close()
				serverURLs = append(serverURLs, s.URL+endpoint)
				servers = append(servers, s)
			}

			c := NewClient()
			resps := c.CollectAllHealthz(serverURLs)

			if len(resps) != len(serverURLs) {
				t.Errorf("CollectAllHealthz should return results from '%v' servers; got '%v'", len(serverURLs), len(resps))
			}

			failedCount := 0
			errorCount := 0
			for _, resp := range resps {
				failedCount = failedCount + resp.Result.Summary.Failed
				if resp.Error != nil {
					errorCount++
				}
			}

			if failedCount != test.failedCount {
				t.Errorf("expected CollectAllHealthz to return '%v' failures; got '%v'", test.failedCount, failedCount)
			}

			if errorCount != test.errorCount {
				t.Errorf("expected CollectAllHealthz to return '%v' errors; got '%v'", test.errorCount, errorCount)
			}
		})
	}
}
