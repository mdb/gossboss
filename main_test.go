package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	exe := exec.Command("go", "build", "-o", "gossboss")
	err := exe.Run()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
