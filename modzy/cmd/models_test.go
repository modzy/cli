package cmd

import (
	"testing"
)

func TestModelsFine(t *testing.T) {
	// for coverage
	runTestCommand(
		[]string{"models"},
		func() {
			Execute()
		},
		nil,
	)
}
