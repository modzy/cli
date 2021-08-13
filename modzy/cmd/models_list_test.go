package cmd

import (
	"net/http"
	"strings"
	"testing"
)

func TestModelsListFine(t *testing.T) {
	x := 0
	stdout, _ := runTestCommand(
		[]string{"models", "list"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected method to be GET, got %s", r.Method)
			}
			if !strings.HasPrefix(r.RequestURI, "/api/models") {
				t.Errorf("get url not expected: %s", r.RequestURI)
			}
			x++
			switch x {
			case 1:
				w.Write([]byte(`[{"modelId": "model1"},{"modelId": "model2"}]`))
			case 2:
				w.Write([]byte(`{"modelId": "model1","modelId":"a1"}`))
			case 3:
				w.Write([]byte(`{"modelId": "model2","modelId":"a2"}`))
			}
		},
	)

	if !strings.Contains(stdout, "model1") {
		t.Fatalf("out not expected: '%s'", stdout)
	}
	if !strings.Contains(stdout, "model2") {
		t.Fatalf("out not expected: '%s'", stdout)
	}
}

func TestModelsListBadFilter(t *testing.T) {
	defer (func() {
		modelsListArgs.Filter = []string{}
	})()

	_, stderr := runTestCommand(
		[]string{"models", "list", "--filter", "b:ad"},
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

func TestModelsListSDKListFailure(t *testing.T) {
	_, stderr := runTestCommand(
		[]string{"models", "list"},
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

func TestModelsListSDKDetailFailure(t *testing.T) {
	x := 0
	_, stderr := runTestCommand(
		[]string{"models", "list"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			x++
			switch x {
			case 1:
				w.Write([]byte(`[{"id": "model1"},{"id": "model2"}]`))
			case 2:
				w.WriteHeader(500)
			}
		},
	)

	if !strings.Contains(stderr, "500") {
		t.Fatalf("out not expected: '%s'", stderr)
	}
}
func TestModelsListOutputerError(t *testing.T) {
	outputer := &modelsListOutputer{}
	err := outputer.Standard(&failWriter{}, []modelSummaryWithMore{})
	if err == nil {
		t.Fatalf("expected an error")
	}
	if err.Error() != "no" {
		t.Errorf("error was not as expected: %v", err)
	}
}
