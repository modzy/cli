package cobrainit

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envPrefix = "MODZY"
)

func InitializeConfig(cmd *cobra.Command, profile string) (string, error) {
	v := viper.New()

	// config files
	// if we have a profile provided, then only read that profile
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/modzy/")
	v.AddConfigPath("$HOME/.modzy/")
	configName := "default"
	if profile != "" {
		configName = profile
	}
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		// if the config file is not found, continue on
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return "", err
		}
	}

	if v.ConfigFileUsed() != "" {
		logrus.Debugf("Using config file at %s", v.ConfigFileUsed())
	}

	// env
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	// flags
	bindFlags(cmd, v)

	return v.ConfigFileUsed(), nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Name == "Profile" {
			// don't do this for the profile flag
			return
		}

		// Environment variables can't have dashes in them, so bind them to their equivalent
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			bindName := fmt.Sprintf("%s_%s", envPrefix, envVarSuffix)
			if err := v.BindEnv(f.Name, bindName); err != nil {
				logrus.WithError(err).Fatalf("Failed to bind flags for %s to %s", f.Name, bindName)
			}
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			if err := cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)); err != nil {
				logrus.WithError(err).Fatalf("Failed to set cmd flags from viper")
			}
		}
	})
}
