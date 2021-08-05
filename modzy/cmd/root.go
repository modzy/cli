package cmd

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	envPrefix = "MODZY"
)

var rootArgs struct {
	Verbose              bool
	Profile              string
	configurationFoundAt string
	BaseURL              string
	APIKey               string
	TeamID               string
	TeamToken            string
}

var rootCmd = &cobra.Command{
	Use:   "modzy",
	Short: "modzy cli",
	Long: `You can provide your authentication token through any mixture of command flags,
ENV variables, or configuration files.  Precedence is command flag > ENV > configuration file.

The flags are:
- base-url
- api-key
- team-id
- team-token

The ENV variables are:
- MODZY_BASE_URL
- MODZY_API_KEY
- MODZY_TEAM_ID
- MODZY_TEAM_TOKEN

To use a configuration file, create a yaml file at any of these locations:
  - /etc/modzy/{profile}.yaml
  - $HOME/.modzy/{profile}.yaml.

The default profile is "default".

This file should look something like:

	> cat ~/.modzy/default.yaml
	base-url: https://base.url
	# use an api key:
	api-key: yourkey.here
	# or use a team key:
	team-id: yourteamid
	team-token: yourteam.token

You can troubleshoot your configuration using the "whoami" command:

	> modzy --profile dev whoami
	Configuration file: /home/user/.modzy/dev.yaml
			Base URL: base
			API Key: yourkey.***
			Team ID: yourteamid
			Team Token: yourteam.***

`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		return initializeConfig(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(configureLogging)

	rootCmd.PersistentFlags().BoolVarP(&rootArgs.Verbose, "verbose", "v", false, "enable debug log output")
	rootCmd.PersistentFlags().StringVarP(&rootArgs.Profile, "profile", "p", "default", "use a profile located at $HOME/.modzy/{profile}")
	rootCmd.PersistentFlags().StringVarP(&rootArgs.BaseURL, "base-url", "", "", "modzy API base URL")
	rootCmd.PersistentFlags().StringVarP(&rootArgs.APIKey, "api-key", "", "", "modzy API key to use for authentication")
	rootCmd.PersistentFlags().StringVarP(&rootArgs.TeamID, "team-id", "", "", "modzy API team ID to use for team authentication")
	rootCmd.PersistentFlags().StringVarP(&rootArgs.TeamToken, "team-token", "", "", "modzy API team token to use for team authentication")
}

// Execute -
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if rootArgs.Verbose {
			logrus.WithError(err).Fatal("Error executing command")
		}
	}
}

func configureLogging() {
	logrus.SetFormatter(&prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.999999999Z07:00",
	})
	if rootArgs.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	// config files
	// if we have a profile provided, then only read that profile
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/modzy/")
	v.AddConfigPath("$HOME/.modzy/")
	configName := "default"
	if rootArgs.Profile != "" {
		configName = rootArgs.Profile
	}
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		// if the config file is not found, continue on
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	rootArgs.configurationFoundAt = v.ConfigFileUsed()
	logrus.Debugf("Using config file at %s", rootArgs.configurationFoundAt)

	// env
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	// flags
	bindFlags(cmd, v)

	return nil
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
