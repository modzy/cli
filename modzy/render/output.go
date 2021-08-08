package render

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/sirupsen/logrus"
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

	if strings.HasPrefix(config, "jsonpath") {
		configSplit := strings.Split(config, "=")
		if len(configSplit) != 2 {
			return fmt.Errorf("jsonpath configuration is invalid")
		}
		jsonpathConfig := strings.TrimRight(strings.TrimLeft(strings.TrimSpace(configSplit[1]), "'"), "'")
		logrus.Debugf("Using jsonpath configuration: '%s'", jsonpathConfig)

		// get a generic unmarshalled json object to query against
		v := interface{}(nil)
		jsonBytes, err := json.Marshal(out)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(jsonBytes, &v); err != nil {
			return err
		}

		res, err := jsonpath.Get(jsonpathConfig, v)
		if err != nil {
			return fmt.Errorf("jsonpath error: %v", err)
		}

		// fmt.Printf("\n\n%v\n\n", res)

		return json.NewEncoder(w).Encode(res)
	}

	return outputer.Standard(os.Stdout, out)
}
