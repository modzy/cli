package cmd

import (
	"os"
	"strings"
	"testing"
)

func TestWhoamiFine(t *testing.T) {
	// for coverage
	os.Args = []string{os.Args[0], "whoami"}
	Execute()
}

func TestWhoamiTeamFine(t *testing.T) {
	// for coverage
	rootArgs.TeamToken = "notsensitive.sensitive"

	out, err := runTestCommand(func() {
		Execute()
	}, []string{"whoami"})

	if err != nil {
		t.Fatalf("error not expected: %v", err)
	}
	if !strings.Contains(out, "Configuration file:") {
		t.Fatalf("out not expected: '%s'", out)
	}
}
