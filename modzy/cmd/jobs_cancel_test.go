package cmd

import (
	"net/http"
	"strings"
	"testing"
)

func TestJobsCancelFine(t *testing.T) {
	runTestCommand(
		[]string{"jobs", "cancel", "theid"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "DELETE" {
				t.Errorf("expected method to be GET, got %s", r.Method)
			}
			if !strings.HasPrefix(r.RequestURI, "/api/jobs/theid") {
				t.Errorf("get url not expected: %s", r.RequestURI)
			}
			w.Write([]byte(`{}`))
		},
	)
}
func TestJobsCancelSDKFailure(t *testing.T) {
	_, stderr := runTestCommand(
		[]string{"jobs", "cancel", "theid"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		},
	)

	if !strings.Contains(stderr, "500") {
		t.Fatalf("out not expected: '%s'", stderr)
	}
}
