package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		arg            string
		expectedOutput string
		expectedErr    string
	}{{
		arg:            "--foo",
		expectedOutput: "Error: unknown flag: --foo\n",
		expectedErr:    "exit status 1",
	}, {
		arg:            "foo",
		expectedOutput: "Error: unknown command \"foo\" for \"gossboss\"\nRun 'gossboss --help' for usage.\n",
		expectedErr:    "exit status 1",
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf("when passed '%s'", test.arg), func(t *testing.T) {
			output, err := exec.Command("./gossboss", test.arg).CombinedOutput()

			if err.Error() != test.expectedErr {
				t.Errorf("expected '%s' to error with '%s'; got '%s'", test.arg, test.expectedErr, err.Error())
			}

			outStr := string(output)
			if outStr != test.expectedOutput {
				t.Errorf("expected '%s' to output '%s'; got '%s'", test.arg, test.expectedOutput, outStr)
			}
		})
	}
}
