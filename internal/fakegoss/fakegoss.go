package fakegoss

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// NewServer creates a new fake goss server for testing.
func NewServer(path, body string, responseCode int) *httptest.Server {
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

// ResponseBody returns a new fake goss server JSON response body for testing.
func ResponseBody(isSuccess bool) string {
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
			"summary-line": "Count: 1, Failed: %v, Duration: 0.048s",
			"test-count": 1,
			"total-duration": 48441102
		}
	}`, isSuccess, failedCount, failedCount)
}
