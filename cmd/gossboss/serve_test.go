package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/mdb/gossboss/internal/fakegoss"
)

func TestServe_WhenNoServerIsStarted(t *testing.T) {
	var (
		description string = "Collect and report goss test results from multiple goss servers' '/healthz' endpoints via a web server JSON endpoint"
	)

	tests := []struct {
		name    string
		arg     string
		outputs []string
		err     error
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
			arg := test.arg
			cmd := exec.Command("./gossboss", "serve", arg)

			stdOut, err := cmd.StdoutPipe()
			if err != nil {
				t.Errorf("expected creation of '%s' command stdout pipe not to error; got '%v'", arg, err)
			}

			stdErr, err := cmd.StderrPipe()
			if err != nil {
				t.Errorf("expected creation of '%s' command stderr pipe not to error; got '%v'", arg, err)
			}

			if err := cmd.Start(); err != nil {
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
				if !strings.Contains(string(stdOutStr), o) {
					t.Errorf("expected 'serve %s' stdout output to include output '%s'; got '%s'", test.arg, o, stdOutStr)
				}
			}
		})
	}
}

func TestServe_WhenServerIsStarted(t *testing.T) {
	placeholderText := "REPLACE_ME"

	tests := []struct {
		name                 string
		arg                  string
		response             *response
		expectedResponseCode int
	}{{
		name:                 "when passed a '--server' with a legitimate server that responds 200",
		arg:                  fmt.Sprintf("--servers=%s", placeholderText),
		expectedResponseCode: 200,
		response: &response{
			code: 200,
			body: fakegoss.ResponseBody(true),
		},
	}, {
		name:                 "when passed a '--server' with a legitimate server that responds 500",
		arg:                  fmt.Sprintf("--servers=%s", placeholderText),
		expectedResponseCode: 500,
		response: &response{
			code: 500,
			body: fakegoss.ResponseBody(false),
		},
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf(test.name), func(t *testing.T) {
			arg := test.arg

			server := fakegoss.NewServer("/healthz", test.response.body, test.response.code)
			defer server.Close()

			arg = strings.ReplaceAll(arg, placeholderText, server.URL+"/healthz")
			cmd := exec.Command("./gossboss", "serve", arg)

			if err := cmd.Start(); err != nil {
				t.Errorf("expected starting '%s' command not to error; got '%v'", arg, err)
			}

			// prevent race condition & ensure server has started
			// TODO: handle the potential for a race condition in a better way
			time.Sleep(1 * time.Second)

			req, err := http.NewRequest("GET", "http://127.0.0.1:8085/healthzs", nil)
			if err != nil {
				t.Errorf("expected request creation when testing '%s' command not to error; got '%s'", arg, err.Error())
			}

			c := &http.Client{}
			resp, err := c.Do(req)
			if err != nil {
				t.Errorf("expected request when testing '%s' command not to error; got '%s'", arg, err.Error())
			}

			if resp.StatusCode != test.expectedResponseCode {
				t.Errorf("expected request to '%s'-started server to respond '%v'; got '%v'", arg, test.expectedResponseCode, resp.StatusCode)
			}

			if err := cmd.Process.Kill(); err != nil {
				t.Errorf("unexpected error killing 'serve %s' proccess: %s", test.arg, err.Error())
			}
		})
	}
}
