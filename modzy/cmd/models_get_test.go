package cmd

import (
	"net/http"
	"strings"
	"testing"

	modzysdkmodel "github.com/modzy/sdk-go/model"
)

func TestModelsGetFine(t *testing.T) {
	stdout, _ := runTestCommand(
		[]string{"models", "get", "theid"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected method to be GET, got %s", r.Method)
			}
			if !strings.HasPrefix(r.RequestURI, "/api/models/theid") {
				t.Errorf("get url not expected: %s", r.RequestURI)
			}
			w.Write([]byte(`{"modelId": "model1"}`))
		},
	)

	if !strings.Contains(stdout, "model1") {
		t.Fatalf("out not expected: '%s'", stdout)
	}
}
func TestModelsGetSDKFailure(t *testing.T) {
	_, stderr := runTestCommand(
		[]string{"models", "get", "theid"},
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

func TestModelsGetOutputerError(t *testing.T) {
	outputer := &modelsGetOutputer{}
	err := outputer.Standard(&failWriter{}, modzysdkmodel.ModelDetails{})
	if err == nil {
		t.Fatalf("expected an error")
	}
	if err.Error() != "no" {
		t.Errorf("error was not as expected: %v", err)
	}
}
