package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestExecuteError(t *testing.T) {
	// for coverage
	rootArgs.Verbose = true
	os.Args = []string{os.Args[0], "--bad"}
	Execute()
}

func TestExecuteFine(t *testing.T) {
	// for coverage
	os.Args = []string{os.Args[0], "whoami"}
	Execute()
}

func TestExecuteRootIsHelp(t *testing.T) {
	// for coverage
	os.Args = []string{os.Args[0]}
	Execute()
}

// helper for running a command with some args
func runTestCommand(callback func(), args []string) (string, error) {
	restoreStdout := os.Stdout
	defer (func() {
		os.Stdout = restoreStdout
	})()

	r, w, _ := os.Pipe()
	defer w.Close()

	os.Stdout = w

	os.Args = append([]string{os.Args[0]}, args...)
	callback()
	w.Close()

	out, err := ioutil.ReadAll(r)
	return string(out), err
}
