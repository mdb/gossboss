package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	tests := []struct {
		arg string
	}{{
		arg: "--help",
	}, {
		arg: "help",
	}, {
		arg: "",
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf("when passed '%s'", test.arg), func(t *testing.T) {
			description := "A tool for collecting goss test results from multiple goss servers' '/healthz' endpoints"

			output, err := exec.Command("./gossboss", test.arg).CombinedOutput()
			if err != nil {
				t.Errorf("expected '%s' not to error; got '%v'", test.arg, err)
			}

			if !strings.Contains(string(output), description) {
				t.Errorf("expected '%s' to output '%s'", test.arg, output)
			}

			if !strings.Contains(string(output), "Usage:") {
				t.Errorf("expected '%s' to output to report usage", test.arg)
			}

			if !strings.Contains(string(output), "Available Commands:") {
				t.Errorf("expected '%s' to output to report available commands", test.arg)
			}
		})
	}
}
