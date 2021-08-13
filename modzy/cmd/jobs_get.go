package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/modzy/cli/internal/render"
	modzysdk "github.com/modzy/sdk-go"
	modzysdkmodel "github.com/modzy/sdk-go/model"

	"github.com/spf13/cobra"
)

var jobsGetArgs struct {
	Output string
}

func init() {
	jobsGetCmd.Flags().StringVarP(&jobsGetArgs.Output, "output", "o", "", "")

	jobsCmd.AddCommand(jobsGetCmd)
}

var jobsGetCmd = &cobra.Command{
	Use:   "get [jobIdentifier]",
	Short: "Get detailed information about a job",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE:  jobsGetRun,
}

func jobsGetRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	jobID := args[0]
	out, err := client.Jobs().GetJobDetails(ctx, &modzysdk.GetJobDetailsInput{
		JobIdentifier: jobID,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	render.Output(os.Stdout, &JobRenderer{}, out.Details, jobsGetArgs.Output)
	return nil
}

type JobRenderer struct{}

func (o *JobRenderer) Standard(w io.Writer, generic interface{}) error {
	out := generic.(modzysdkmodel.JobDetails)

	tabbed := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(tabbed, "ID: \t%s\n", out.JobIdentifier)
	fmt.Fprintf(tabbed, "Status: \t%s\n", out.Status)
	fmt.Fprintf(tabbed, "Model: \t%s\n", out.Model.Name)
	fmt.Fprintf(tabbed, "       \t%s@%s\n", out.Model.Identifier, out.Model.Version)

	if err := tabbed.Flush(); err != nil {
		return err
	}

	return nil
}
