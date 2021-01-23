package gossboss_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdb/gossboss"
	"github.com/mdb/gossboss/internal/fakegoss"
)

func TestHandleHealthzs(t *testing.T) {
	tests := []struct {
		name         string
		expectedCode int
		response     testResponse
	}{{
		name:         "the backend Goss servers responds 200 and there are no failures",
		expectedCode: http.StatusOK,
		response: testResponse{
			code: 200,
			body: fakegoss.ResponseBody(true),
		}}, {
		name:         "a backend Goss server responds 200 but with invalid JSON",
		expectedCode: http.StatusInternalServerError,
		response: testResponse{
			code: 200,
			body: "foo",
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			endpoint := "/healthz"
			gossServer := fakegoss.NewServer(endpoint, test.response.body, test.response.code)
			defer gossServer.Close()

			req, err := http.NewRequest("GET", "/healthzs", nil)
			if err != nil {
				t.Fatal(err)
			}
			server := gossboss.NewServer(":8081", []string{gossServer.URL + endpoint})
			handler := http.HandlerFunc(server.HandleHealthzs)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != test.expectedCode {
				t.Errorf("expected HandleHealthzs to return '%v'; got '%v'", test.expectedCode, status)
			}
		})
	}
}
