package cmd

import (
	"net/http"
	"strings"
	"testing"

	modzysdkmodel "github.com/modzy/sdk-go/model"
)

func TestProjectsListFine(t *testing.T) {
	stdout, stderr := runTestCommand(
		[]string{"projects", "list"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected method to be GET, got %s", r.Method)
			}
			if !strings.HasPrefix(r.RequestURI, "/api/accounting/projects") {
				t.Errorf("get url not expected: %s", r.RequestURI)
			}
			w.Write([]byte(`[{"identifier": "project1","accessKeys":[{"prefix":"theprefix"}]},{"identifier": "project2"}]`))
		},
	)

	if !strings.Contains(stdout, "project1") {
		t.Fatalf("out not expected: '%s' | '%s'", stdout, stderr)
	}
	if !strings.Contains(stdout, "project2") {
		t.Fatalf("out not expected: '%s'", stdout)
	}
}

func TestProjectsListBadFilter(t *testing.T) {
	defer (func() {
		projectsListArgs.Filter = []string{}
	})()

	_, stderr := runTestCommand(
		[]string{"projects", "list", "--filter", "b:ad"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected method to be GET, got %s", r.Method)
			}
			if r.RequestURI != "/api/accounts/projects" {
				t.Errorf("get url not expected: %s", r.RequestURI)
			}
			w.Write([]byte(`[{"identifier": "jsonID"}]`))
		},
	)

	if !strings.Contains(stderr, "Filter is not correctly formatted") {
		t.Fatalf("out not expected: '%s'", stderr)
	}
}

func TestProjectsListSDKListFailure(t *testing.T) {
	stdout, stderr := runTestCommand(
		[]string{"projects", "list"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		},
	)

	if !strings.Contains(stderr, "500") {
		t.Fatalf("out not expected: '%s' | '%s'", stdout, stderr)
	}
}

func TestProjectsListOutputerError(t *testing.T) {
	outputer := &projectsListOutputer{}
	err := outputer.Standard(&failWriter{}, []modzysdkmodel.AccountingProject{})
	if err == nil {
		t.Fatalf("expected an error")
	}
	if err.Error() != "no" {
		t.Errorf("error was not as expected: %v", err)
	}
}
