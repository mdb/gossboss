package gossboss_test

import (
	"net/http/httptest"
	"testing"

	"github.com/mdb/gossboss"
	"github.com/mdb/gossboss/internal/fakegoss"
)

func BenchmarkCollectHealthzs(b *testing.B) {
	servers := []*httptest.Server{}
	serverURLs := []string{}
	endpoint := "/healthz"

	for i := 1; i < 5; i++ {
		s := fakegoss.NewServer(endpoint, fakegoss.ResponseBody(true), 200)
		defer s.Close()

		serverURLs = append(serverURLs, s.URL+endpoint)
		servers = append(servers, s)
	}

	c := gossboss.NewClient()

	for i := 0; i < b.N; i++ {
		_ = c.CollectHealthzs(serverURLs)
	}
}
