package cmd

import (
	"net/http"
	"strings"
	"testing"
)

func TestAlertsFine(t *testing.T) {
	x := 0
	stdout, stderr := runTestCommand(
		[]string{"alerts"},
		func() {
			Execute()
		},
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("expected method to be GET, got %s", r.Method)
			}
			x++
			switch x {
			case 1:
				if !strings.HasPrefix(r.RequestURI, "/api/notifications/alerts") {
					t.Errorf("get url not expected: %s", r.RequestURI)
				}
				w.Write([]byte(`[{"type": "type1", "count": 10},{"type": "type2", "count": 20}]`))
			case 2:
				if !strings.HasPrefix(r.RequestURI, "/api/notifications/alerts/type1") {
					t.Errorf("get url not expected: %s", r.RequestURI)
				}
				w.Write([]byte(`["type1-ent1", "type1-ent2"]`))
			case 3:
				w.Write([]byte(`["type2-ent1", "type2-ent2"]`))
			}
		},
	)

	if !strings.Contains(stdout, "type1") {
		t.Fatalf("out not expected: '%s' | %s", stdout, stderr)
	}
	if !strings.Contains(stdout, "type2") {
		t.Fatalf("out not expected: '%s' | %s", stdout, stderr)
	}
}

func TestAlertsSDKListFailure(t *testing.T) {
	_, stderr := runTestCommand(
		[]string{"alerts"},
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

func TestAlertsSDKDetailFailure(t *testing.T) {
	x := 0
	_, stderr := runTestCommand(
		[]string{"alerts"},
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
func TestAlertsOutputerError(t *testing.T) {
	outputer := &alertsGetOutputer{}
	err := outputer.Standard(&failWriter{}, []alertAndDetail{})
	if err == nil {
		t.Fatalf("expected an error")
	}
	if err.Error() != "no" {
		t.Errorf("error was not as expected: %v", err)
	}
}
