package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/mdb/gossboss/internal/fakegoss"
)

var (
	errExit     error  = errors.New("exit status 1")
	description string = "Collect and report goss test results from multiple goss servers' '/healthz' endpoints"
	replaceText string = "REPLACE_ME"
)

func TestCollect(t *testing.T) {
	type response struct {
		body string
		code int
	}

	tests := []struct {
		arg      string
		outputs  []string
		err      error
		response *response
	}{{
		arg: "--help",
		outputs: []string{
			description,
			"Usage:",
			"Flags:",
		},
		err: nil,
	}, {
		arg: "-h",
		outputs: []string{
			description,
			"Usage:",
			"Flags:",
		},
		err: nil,
	}, {
		arg: "",
		outputs: []string{
			"Error: required flag(s) \"servers\" not set",
		},
		err: errExit,
	}, {
		arg: "--servers",
		outputs: []string{
			"Error: flag needs an argument: --servers",
		},
		err: errExit,
	}, {
		arg: fmt.Sprintf("--servers=%s", replaceText),
		outputs: []string{
			"Error: invalid character 'o' in literal false (expecting 'a')",
			"Goss test collection error",
		},
		err: errExit,
		response: &response{
			code: 200,
			body: "foo",
		},
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf("when 'collect' is passed '%s'", test.arg), func(t *testing.T) {
			arg := test.arg

			if test.response != nil {
				server := fakegoss.NewServer("/healthz", test.response.body, test.response.code)
				arg = strings.ReplaceAll(arg, replaceText, server.URL+"/healthz")
				defer server.Close()
			}

			output, err := exec.Command("./gossboss", "collect", arg).CombinedOutput()

			if test.err == nil && err != nil {
				t.Errorf("expected '%s' not to error; got '%v'", arg, err)
			}

			if test.err != nil && err == nil {
				t.Errorf("expected '%s' to error with '%s', but it didn't error", arg, test.err.Error())
			}

			if test.err != nil && err != nil && test.err.Error() != err.Error() {
				t.Errorf("expected '%s' to error with '%s'; got '%s'", arg, test.err.Error(), err.Error())
			}

			for _, o := range test.outputs {
				if !strings.Contains(string(output), o) {
					t.Errorf("expected '%s' to include output '%s'; got '%s'", test.arg, o, output)
				}
			}
		})
	}
}
