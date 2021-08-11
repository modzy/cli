package render

import (
	"encoding/json"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type StandardOutputer interface {
	Standard(w io.Writer, generic interface{}) error
}

func Output(w io.Writer, outputer StandardOutputer, out interface{}, config string) error {
	if config == "json" {
		return json.NewEncoder(w).Encode(out)
	}

	if config == "yaml" {
		return yaml.NewEncoder(w).Encode(out)
	}

	return outputer.Standard(os.Stdout, out)
}
