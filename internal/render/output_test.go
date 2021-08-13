package render_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/modzy/cli/internal/render"
)

var testOut = map[string]string{
	"a": "b",
}

type nojson struct{}

func (no nojson) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("no json")
}

type noyaml struct{}

func (no noyaml) MarshalYAML() (interface{}, error) {
	return nil, fmt.Errorf("no yaml")
}

func TestOutputJson(t *testing.T) {
	var buf bytes.Buffer

	// can't marshal
	err := render.Output(&buf, nil, nojson{}, "json")
	if err == nil {
		t.Fatalf("Err was expected")
	}

	// good
	err = render.Output(&buf, nil, testOut, "json")
	if err != nil {
		t.Fatalf("Err not expected: %v", err)
	}
	out, err := ioutil.ReadAll(&buf)
	if err != nil {
		t.Errorf("Err not expected: %v", err)
	}
	if strings.TrimSpace(string(out)) != `{"a":"b"}` {
		t.Errorf("json out not correct: '%s'", out)
	}
}

func TestOutputYaml(t *testing.T) {
	var buf bytes.Buffer

	// can't marshal
	err := render.Output(&buf, nil, noyaml{}, "yaml")
	if err == nil {
		t.Fatalf("Err was expected")
	}

	// good
	err = render.Output(&buf, nil, testOut, "yaml")
	if err != nil {
		t.Fatalf("Err not expected: %v", err)
	}
	out, err := ioutil.ReadAll(&buf)
	if err != nil {
		t.Errorf("Err not expected: %v", err)
	}
	if strings.TrimSpace(string(out)) != `a: b` {
		t.Errorf("json out not correct: '%s'", out)
	}
}

func TestOutputJsonPath(t *testing.T) {
	var buf bytes.Buffer

	err := render.Output(&buf, nil, nil, "jsonpath-junk")
	if err == nil {
		t.Fatalf("Err was expected")
	}
	if !strings.Contains(err.Error(), "jsonpath configuration is") {
		t.Errorf("Err was not expected kind: %v", err)
	}

	// bad json
	err = render.Output(&buf, nil, nojson{}, "jsonpath=donotcare")
	if err == nil {
		t.Fatalf("Err was expected")
	}

	// bad jsonpath
	err = render.Output(&buf, nil, testOut, "jsonpath=;;")
	if err == nil {
		t.Fatalf("Err was expected")
	}
	if !strings.Contains(err.Error(), "jsonpath error:") {
		t.Errorf("Err was not expected kind: %v", err)
	}

	// good
	err = render.Output(&buf, nil, testOut, "jsonpath=$.a")
	if err != nil {
		t.Fatalf("Err was not expected: %v", err)
	}
	out, err := ioutil.ReadAll(&buf)
	if err != nil {
		t.Errorf("Err not expected: %v", err)
	}
	if strings.TrimSpace(string(out)) != `"b"` {
		t.Errorf("standard out not correct: '%s'", out)
	}
}

type testOutputer struct{}

func (o *testOutputer) Standard(w io.Writer, generic interface{}) error {
	if generic == nil {
		return fmt.Errorf("no generic")
	}
	fmt.Fprintf(w, "written")
	return nil
}

func TestOutputStandard(t *testing.T) {
	var buf bytes.Buffer

	// outputer failure
	err := render.Output(&buf, &testOutputer{}, nil, "anything-else")
	if err == nil {
		t.Fatalf("Err was expected")
	}

	// good
	err = render.Output(&buf, &testOutputer{}, testOut, "anything-else")
	if err != nil {
		t.Fatalf("Err not expected: %v", err)
	}
	out, err := ioutil.ReadAll(&buf)
	if err != nil {
		t.Errorf("Err not expected: %v", err)
	}
	if strings.TrimSpace(string(out)) != `written` {
		t.Errorf("standard out not correct: '%s'", out)
	}
}
