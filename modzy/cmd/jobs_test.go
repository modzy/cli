package cmd

import (
	"os"
	"testing"
)

func TestJobsFine(t *testing.T) {
	// for coverage
	os.Args = []string{os.Args[0], "jobs"}
	Execute()
}
