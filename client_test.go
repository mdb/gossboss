package gossboss

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/mdb/gossboss/internal/fakegoss"
)

type testResponse struct {
	code int
	body string
}

func TestGetHealthz(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
		failedCount int
		response    testResponse
	}{{
		name:        "the server responds 200 and there are no failures",
		expectedErr: nil,
		failedCount: 0,
		response: testResponse{
			code: 200,
			body: fakegoss.ResponseBody(true),
		}}, {
		name:        "the server responds 500 and there are failures",
		expectedErr: nil,
		failedCount: 1,
		response: testResponse{
			code: 500,
			body: fakegoss.ResponseBody(false),
		}}, {
		name:        "the server responds 200, but serves invalid JSON",
		expectedErr: errors.New("invalid character 'o' in literal false (expecting 'a')"),
		failedCount: 0,
		response: testResponse{
			code: 200,
			body: "foo",
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			endpoint := "/healthz"
			server := fakegoss.NewServer(endpoint, test.response.body, test.response.code)
			defer server.Close()

			c := NewClient()
			resp, err := c.GetHealthz(server.URL + endpoint)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected GetHealthz to return error '%v'; got '%v'", test.expectedErr, err)
			}

			if err != nil && resp != nil {
				t.Errorf("expected GetHealthz to return a nil response given a non-nil error of '%v'", err)
			}

			if resp != nil && resp.Summary.Failed != test.failedCount {
				t.Errorf("expected GetHealthz to return a Summary.Failed of '%v'; got '%v'", test.failedCount, resp.Summary.Failed)
			}
		})
	}
}

func TestCollectAllHealthz(t *testing.T) {
	tests := []struct {
		name        string
		failedCount int
		errorCount  int
		responses   []testResponse
	}{{
		name:        "all servers respond 200",
		failedCount: 0,
		errorCount:  0,
		responses: []testResponse{{
			code: 200,
			body: fakegoss.ResponseBody(true),
		}, {
			code: 200,
			body: fakegoss.ResponseBody(true),
		}}}, {
		name:        "1 server responds 500",
		failedCount: 1,
		errorCount:  0,
		responses: []testResponse{{
			code: 500,
			body: fakegoss.ResponseBody(false),
		}, {
			code: 200,
			body: fakegoss.ResponseBody(true),
		}}}, {
		name:        "1 server returns invalid JSON",
		failedCount: 0,
		errorCount:  1,
		responses: []testResponse{{
			code: 200,
			body: "foo",
		}, {
			code: 200,
			body: fakegoss.ResponseBody(true),
		}},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			endpoint := "/healthz"

			servers := []*httptest.Server{}
			serverURLs := []string{}
			for _, r := range test.responses {
				s := fakegoss.NewServer(endpoint, r.body, r.code)
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
				if resp.Result != nil {
					failedCount = failedCount + resp.Result.Summary.Failed
				}

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
