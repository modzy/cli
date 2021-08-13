package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestExecuteError(t *testing.T) {
	_, stderr := runTestCommand(
		[]string{"--verbose", "--bad"},
		func() {
			Execute()
		},
		nil,
	)

	if !strings.Contains(stderr, "unknown flag: --bad") {
		t.Fatalf("out not expected: '%s'", stderr)
	}
}

func TestExecuteFine(t *testing.T) {
	// for coverage
	runTestCommand(
		[]string{"--verbose", "whoami"},
		func() {
			Execute()
		},
		nil,
	)
}

func TestExecuteRootIsHelp(t *testing.T) {
	// for coverage
	runTestCommand(
		[]string{},
		func() {
			Execute()
		},
		nil,
	)
}

// helper for running a command with some args
func runTestCommand(args []string, callback func(), handler func(w http.ResponseWriter, r *http.Request)) (string, string) {
	defer (func() {
		// persistent flag resets that we care to not have pass through test methods
		rootArgs.Verbose = false
		rootArgs.VerboseHTTP = false
	})()

	if handler == nil {
		handler = func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{}`))
		}
	}

	serv := httptest.NewServer(http.HandlerFunc(handler))
	defer serv.Close()

	restoreStdout := os.Stdout
	restoreStderr := os.Stderr

	defer (func() {
		os.Stdout = restoreStdout
		os.Stderr = restoreStderr
	})()

	r, w, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()

	os.Stdout = w
	os.Stderr = wErr

	os.Args = append([]string{"fake-executable", "--base-url", serv.URL}, args...)

	callback()

	w.Close()
	wErr.Close()

	stdinContent, _ := ioutil.ReadAll(r)
	stderrContent, _ := ioutil.ReadAll(rErr)

	return string(stdinContent), string(stderrContent)
}

type failWriter struct{}

func (fw *failWriter) Write(bytes []byte) (int, error) {
	return 0, fmt.Errorf("no")
}
