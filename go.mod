module github.com/modzy/cli

go 1.16

require (
	github.com/PaesslerAG/jsonpath v0.1.1
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/modzy/sdk-go v0.0.5
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/modzy/sdk-go v0.0.5 => ../sdk-go
