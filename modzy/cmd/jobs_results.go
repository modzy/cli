package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/modzy/cli/internal/render"
	modzysdk "github.com/modzy/sdk-go"
	modzysdkmodel "github.com/modzy/sdk-go/model"

	"github.com/spf13/cobra"
)

var jobsResultsArgs struct {
	Output string
}

func init() {
	jobsResultsCmd.Flags().StringVarP(&jobsResultsArgs.Output, "output", "o", "", "")

	jobsCmd.AddCommand(jobsResultsCmd)
}

var jobsResultsCmd = &cobra.Command{
	Use:          "results [jobIdentifier]",
	Short:        "Get job results",
	Long:         ``,
	Args:         cobra.ExactArgs(1),
	RunE:         jobsResultsRun,
	SilenceUsage: true,
}

func jobsResultsRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	jobID := args[0]
	out, err := client.Jobs().GetJobResults(ctx, &modzysdk.GetJobResultsInput{
		JobIdentifier: jobID,
	})
	if err != nil {
		return err
	}

	render.Output(os.Stdout, &jobsResultsRenderer{}, out.Results, jobsResultsArgs.Output)
	return nil
}

type jobsResultsRenderer struct{}

func (o *jobsResultsRenderer) Standard(w io.Writer, generic interface{}) error {
	out := generic.(modzysdkmodel.JobResults)

	totalProcessing := 0
	for _, res := range out.Results {
		totalProcessing += res.ElapsedTime
	}
	for _, res := range out.Failures {
		totalProcessing += res.ElapsedTime
	}

	tabbed := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(tabbed, "ID: \t%s\n", out.JobIdentifier)
	fmt.Fprintf(tabbed, "Finished: \t%t (%d Total, %d Completed, %d Failed)\n", out.Finished, out.Total, out.Completed, out.Failed)
	fmt.Fprintf(tabbed, "Start: \t%s\n", out.SubmittedAt)
	fmt.Fprintf(tabbed, "Processing Time: \t%s\n", time.Duration(totalProcessing)*time.Millisecond)
	if err := tabbed.Flush(); err != nil {
		// for coverage
		return err
	}

	fmt.Fprintf(w, "\nResults:\n")
	tabbed = tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
	fmt.Fprintf(tabbed, "Input\tStatus\tStart Time\tElapsed Time\n")
	fmt.Fprintf(tabbed, "-----\t------\t----------\t------------\n")
	for inputName, res := range out.Results {
		fmt.Fprintf(tabbed, "%s\t%s\t%s\t%s\n",
			inputName, res.Status, res.StartTime, time.Duration(res.ElapsedTime)*time.Millisecond,
		)
	}
	for inputName, res := range out.Failures {
		fmt.Fprintf(tabbed, "%s\t%s\t%s\t%s\n",
			inputName, res.Status, res.StartTime, time.Duration(res.ElapsedTime)*time.Millisecond,
		)
	}
	tabbed.Flush()

	return nil
}
