package cmd

import (
	"net/http"
	"strings"
	"testing"

	modzysdkmodel "github.com/modzy/sdk-go/model"
)

func TestJobsListFine(t *testing.T) {
	stdout, _ := runTestCommand(
		[]string{"jobs", "list"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected method to be GET, got %s", r.Method)
			}
			if !strings.HasPrefix(r.RequestURI, "/api/jobs/history?direction=DESC") {
				t.Errorf("get url not expected: %s", r.RequestURI)
			}
			w.Write([]byte(`[{"jobIdentifier": "job1"},{"jobIdentifier": "job2"}]`))
		},
	)

	if !strings.Contains(stdout, "job1") {
		t.Fatalf("out not expected: '%s'", stdout)
	}
	if !strings.Contains(stdout, "job2") {
		t.Fatalf("out not expected: '%s'", stdout)
	}
}

func TestJobsListBadFilter(t *testing.T) {
	defer (func() {
		jobsListArgs.Filter = []string{}
	})()

	_, stderr := runTestCommand(
		[]string{"jobs", "list", "--filter", "b:ad"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected method to be GET, got %s", r.Method)
			}
			if r.RequestURI != "/api/jobs" {
				t.Errorf("get url not expected: %s", r.RequestURI)
			}
			w.Write([]byte(`[{"jobIdentifier": "jsonID"}]`))
		},
	)

	if !strings.Contains(stderr, "Filter is not correctly formatted") {
		t.Fatalf("out not expected: '%s'", stderr)
	}
}

func TestJobsListSDKFailure(t *testing.T) {
	_, stderr := runTestCommand(
		[]string{"jobs", "list"},
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

func TestJobsListOutputerError(t *testing.T) {
	outputer := &jobsListOutputer{}
	err := outputer.Standard(&failWriter{}, []modzysdkmodel.JobDetails{})
	if err == nil {
		t.Fatalf("expected an error")
	}
	if err.Error() != "no" {
		t.Errorf("error was not as expected: %v", err)
	}
}
