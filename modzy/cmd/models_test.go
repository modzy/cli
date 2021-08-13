package cmd

import (
	"os"
	"testing"
)

func TestModelsFine(t *testing.T) {
	// for coverage
	os.Args = []string{os.Args[0], "models"}
	Execute()
}
