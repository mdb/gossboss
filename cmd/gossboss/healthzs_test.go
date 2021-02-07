package main

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/mdb/gossboss/internal/fakegoss"
)

func TestHealthzs(t *testing.T) {
	var (
		errExit         error  = errors.New("exit status 1")
		description     string = "Collect and report goss test results from multiple goss servers' '/healthz' endpoints"
		placeholderText string = "REPLACE_ME"
	)

	tests := []struct {
		name     string
		arg      string
		outputs  []string
		err      error
		response *response
	}{{
		name: "when passed '--help'",
		arg:  "--help",
		outputs: []string{
			description,
			"Usage:",
			"Flags:",
		},
		err: nil,
	}, {
		name: "when passed '-h'",
		arg:  "-h",
		outputs: []string{
			description,
			"Usage:",
			"Flags:",
		},
		err: nil,
	}, {
		name: "when passed nothing",
		arg:  "",
		outputs: []string{
			"Error: required flag(s) \"servers\" not set",
		},
		err: errExit,
	}, {
		name: "when passed a '--server' with no value",
		arg:  "--servers",
		outputs: []string{
			"Error: flag needs an argument: --servers",
		},
		err: errExit,
	}, {
		name: "when passed a '--server' server that responds 200 but returns invalid JSON",
		arg:  fmt.Sprintf("--servers=%s", placeholderText),
		outputs: []string{
			fmt.Sprintf("✘ %s", placeholderText),
			"Error: invalid character 'o' in literal false (expecting 'a')",
			"Goss test collection error",
		},
		err: errExit,
		response: &response{
			code: 200,
			body: "foo",
		},
	}, {
		name: "when passed a '--server' server that has failing tests",
		arg:  fmt.Sprintf("--servers=%s", placeholderText),
		outputs: []string{
			fmt.Sprintf("✘ %s", placeholderText),
			"Error: Goss test failed",
		},
		err: errExit,
		response: &response{
			code: 500,
			body: fakegoss.ResponseBody(false),
		},
	}, {
		name: "when passed a '--server' server that has no failing tests",
		arg:  fmt.Sprintf("--servers=%s", placeholderText),
		outputs: []string{
			fmt.Sprintf("✔ %s", placeholderText),
		},
		err: nil,
		response: &response{
			code: 200,
			body: fakegoss.ResponseBody(true),
		},
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf(test.name), func(t *testing.T) {
			var server *httptest.Server
			arg := test.arg

			if test.response != nil {
				server = fakegoss.NewServer("/healthz", test.response.body, test.response.code)
				arg = strings.ReplaceAll(arg, placeholderText, server.URL+"/healthz")
				defer server.Close()
			}

			output, err := exec.Command("./gossboss", "healthzs", arg).CombinedOutput()

			if test.err == nil && err != nil {
				t.Errorf("expected 'healthzs %s' not to error; got '%v'", arg, err)
			}

			if test.err != nil && err == nil {
				t.Errorf("expected 'healthzs %s' to error with '%s', but it didn't error", arg, test.err.Error())
			}

			if test.err != nil && err != nil && test.err.Error() != err.Error() {
				t.Errorf("expected 'healthzs %s' to error with '%s'; got '%s'", arg, test.err.Error(), err.Error())
			}

			for _, o := range test.outputs {
				if strings.Contains(o, placeholderText) {
					o = strings.ReplaceAll(o, placeholderText, server.URL+"/healthz")
				}

				if !strings.Contains(string(output), o) {
					t.Errorf("expected 'healthzs %s' to include output '%s'; got '%s'", test.arg, o, output)
				}
			}
		})
	}
}
