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

func TestGetResults(t *testing.T) {
	endpoint := "/healthz"
	respStr := "body"
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
