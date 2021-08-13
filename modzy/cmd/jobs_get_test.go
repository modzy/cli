package cmd

import (
	"net/http"
	"strings"
	"testing"

	modzysdkmodel "github.com/modzy/sdk-go/model"
)

func TestJobsGetFine(t *testing.T) {
	stdout, _ := runTestCommand(
		[]string{"jobs", "get", "theid"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected method to be GET, got %s", r.Method)
			}
			if !strings.HasPrefix(r.RequestURI, "/api/jobs/theid") {
				t.Errorf("get url not expected: %s", r.RequestURI)
			}
			w.Write([]byte(`{"jobIdentifier": "job1"}`))
		},
	)

	if !strings.Contains(stdout, "job1") {
		t.Fatalf("out not expected: '%s'", stdout)
	}
}
func TestJobsGetSDKFailure(t *testing.T) {
	_, stderr := runTestCommand(
		[]string{"jobs", "get", "theid"},
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

func TestJobsGetOutputerError(t *testing.T) {
	outputer := &jobsGetRenderer{}
	err := outputer.Standard(&failWriter{}, modzysdkmodel.JobDetails{})
	if err == nil {
		t.Fatalf("expected an error")
	}
	if err.Error() != "no" {
		t.Errorf("error was not as expected: %v", err)
	}
}
