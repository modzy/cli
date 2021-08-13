package cmd

import (
	"fmt"
	"os"

	modzysdk "github.com/modzy/sdk-go"

	"github.com/spf13/cobra"
)

// var jobsCancelArgs struct {
// }

func init() {
	jobsCmd.AddCommand(jobsCancelCmd)
}

var jobsCancelCmd = &cobra.Command{
	Use:   "cancel [jobIdentifier]",
	Short: "Cancel a job",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE:  jobsCancelRun,
}

func jobsCancelRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	jobID := args[0]
	_, err := client.Jobs().CancelJob(ctx, &modzysdk.CancelJobInput{
		JobIdentifier: jobID,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return nil
}
