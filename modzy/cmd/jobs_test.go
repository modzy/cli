package cmd

import (
	"testing"
)

func TestJobsFine(t *testing.T) {
	// for coverage
	runTestCommand(
		[]string{"jobs"},
		func() {
			Execute()
		},
		nil,
	)
}
