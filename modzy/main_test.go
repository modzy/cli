package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	restoreStdout := os.Stdout
	defer (func() {
		os.Stdout = restoreStdout
	})()

	r, w, _ := os.Pipe()
	defer w.Close()

	os.Stdout = w

	os.Args = []string{os.Args[0], "--help"}
	main()
	w.Close()

	out, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("Did not expect an error: %v", err)
	}
	if !strings.HasPrefix(string(out), "You can provide your authentication token") {
		t.Errorf("out was not as expected: '%s'", out)
	}
}
