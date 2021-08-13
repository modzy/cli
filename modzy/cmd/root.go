package cmd

import (
	"github.com/modzy/cli/internal/cobrainit"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var rootArgs struct {
	Verbose              bool
	VerboseHTTP          bool
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
		profileFound, err := cobrainit.InitializeConfig(cmd, rootArgs.Profile)
		rootArgs.configurationFoundAt = profileFound
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(configureLogging)
	rootCmd.PersistentFlags().BoolVarP(&rootArgs.Verbose, "verbose", "v", false, "enable more verbose log output for debugging purposes")
	rootCmd.PersistentFlags().BoolVarP(&rootArgs.VerboseHTTP, "verbose-http", "", false, "enable log output of http request and response data")
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
			logrus.WithError(err).Error("Error executing command")
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
