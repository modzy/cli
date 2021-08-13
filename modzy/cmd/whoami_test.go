package cmd

import (
	"strings"
	"testing"
)

func TestWhoamiFine(t *testing.T) {
	// for coverage
	runTestCommand(
		[]string{"whoami"},
		func() {
			Execute()
		},
		nil,
	)
}

func TestWhoamiTeamFine(t *testing.T) {
	// for coverage
	stdout, _ := runTestCommand(
		[]string{"--team-token", "notsensitive.sensitive", "--team-id", "teamid", "whoami"},
		func() {
			Execute()
		},
		nil,
	)

	if !strings.Contains(stdout, "Configuration file:") {
		t.Fatalf("out not expected: '%s'", stdout)
	}
}
