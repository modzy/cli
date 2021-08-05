package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var whoamiArgs struct {
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Summarize effective authentication information",
	Long:  `This will get effective authentication information that this command is using and will display the non-sensitive portion for the purpose of troubleshooting authentication.`,
	RunE:  whoamiRun,
}

func whoamiRun(cmd *cobra.Command, args []string) error {
	apiKeySplit := strings.Split(rootArgs.APIKey, ".")
	apiKey := ""
	if len(apiKeySplit) > 1 {
		apiKey = fmt.Sprintf("%s.***", apiKeySplit[0])
	}

	teamTokenSplit := strings.Split(rootArgs.TeamToken, ".")
	teamToken := ""
	if len(teamTokenSplit) > 1 {
		teamToken = fmt.Sprintf("%s.***", teamTokenSplit[0])
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "Configuration file: \t%s\n", rootArgs.configurationFoundAt)
	fmt.Fprintf(w, "Base URL: \t%s\n", rootArgs.BaseURL)
	fmt.Fprintf(w, "API Key: \t%s\n", apiKey)
	fmt.Fprintf(w, "Team ID: \t%s\n", rootArgs.TeamID)
	fmt.Fprintf(w, "Team Token: \t%s\n", teamToken)
	w.Flush()

	return nil
}
