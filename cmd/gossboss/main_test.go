package main

import (
	"os"
	"os/exec"
	"testing"
)

// response can be used to express a desired fakegoss.NewServer
// HTTP response body and response code.
type response struct {
	body string
	code int
}

func TestMain(m *testing.M) {
	// compile a 'gossboss' for for use in running tests
	exe := exec.Command("go", "build", "-ldflags", "-X main.version=test", "-o", "gossboss")
	err := exe.Run()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
