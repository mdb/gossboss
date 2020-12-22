package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	// compile a 'gossboss' for for use in running tests
	exe := exec.Command("go", "build", "-o", "gossboss")
	err := exe.Run()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
