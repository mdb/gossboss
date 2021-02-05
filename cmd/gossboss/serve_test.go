package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/mdb/gossboss/internal/fakegoss"
)

func TestServe(t *testing.T) {
	var (
		description     string = "Collect and report goss test results from multiple goss servers' '/healthz' endpoints via a web server JSON endpoint"
		placeholderText string = "REPLACE_ME"
	)

	type response struct {
		body string
		code int
	}

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
		name:    "when passed nothing",
		arg:     "",
		outputs: []string{},
		err:     errors.New("Error: required flag(s) \"servers\" not set"),
	}, {
		name:    "when passed a '--server' with no value",
		arg:     "--servers",
		outputs: []string{},
		err:     errors.New("Error: flag needs an argument: --servers"),
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

			cmd := exec.Command("./gossboss", "serve", arg)

			stdOut, err := cmd.StdoutPipe()
			if err != nil {
				t.Errorf("expected creation of '%s' command stdout pipe not to error; got '%v'", arg, err)
			}

			stdErr, err := cmd.StderrPipe()
			if err != nil {
				t.Errorf("expected creation of '%s' command stderr pipe not to error; got '%v'", arg, err)
			}

			err = cmd.Start()
			if err != nil {
				t.Errorf("expected starting '%s' command not to error; got '%v'", arg, err)
			}

			outBuf := new(bytes.Buffer)
			outBuf.ReadFrom(stdOut)
			stdOutStr := outBuf.String()

			errBuf := new(bytes.Buffer)
			errBuf.ReadFrom(stdErr)
			stdErrStr := errBuf.String()

			if test.err == nil && stdErrStr != "" {
				t.Errorf("expected '%s' not to output to stderr, but it output '%s'", arg, stdErrStr)
			}

			if test.err != nil && stdErrStr == "" {
				t.Errorf("expected 'serve %s' stderr output to include '%s', but it didn't write to stderr", arg, test.err.Error())
			}

			if test.err != nil && stdErrStr != "" && !strings.Contains(stdErrStr, test.err.Error()) {
				t.Errorf("expected 'serve %s' stderr output to include '%s'; got '%s'", arg, test.err.Error(), stdErrStr)
			}

			for _, o := range test.outputs {
				if strings.Contains(o, placeholderText) {
					o = strings.ReplaceAll(o, placeholderText, server.URL+"/healthz")
				}

				if !strings.Contains(string(stdOutStr), o) {
					t.Errorf("expected 'serve %s' stdout output to include output '%s'; got '%s'", test.arg, o, stdOutStr)
				}
			}
		})
	}
}
