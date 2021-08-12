package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/modzy/cli/modzy/render"

	modzysdk "github.com/modzy/sdk-go"
	modzysdkmodel "github.com/modzy/sdk-go/model"

	"github.com/spf13/cobra"
)

var jobsListArgs struct {
	pagingArgs
	Output string
}

func init() {
	jobsListCmd.Flags().StringArrayVarP(
		&jobsListArgs.Filter, "filter", "",
		[]string{
			"startDate=T-30",
			"endDate=T",
		},
		"",
	)
	jobsListCmd.Flags().IntVarP(&jobsListArgs.Take, "take", "", 10, "")
	jobsListCmd.Flags().IntVarP(&jobsListArgs.Page, "page", "", 1, "")
	jobsListCmd.Flags().StringVarP(&jobsListArgs.Output, "output", "o", "", "")
	jobsListCmd.Flags().StringVarP(
		&jobsListArgs.Sort, "sort", "",
		string(modzysdk.ListJobsHistorySortFieldCreatedAt)+":desc",
		"",
	)

	jobsCmd.AddCommand(jobsListCmd)
}

var jobsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List jobs",
	Long:  ``,
	RunE:  jobsListRun,
}

func jobsListRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	paging, err := jobsListArgs.GetPagingInput()
	if err != nil {
		return err
	}

	input := &modzysdk.ListJobsHistoryInput{
		Paging: paging,
	}
	out, err := client.Jobs().ListJobsHistory(ctx, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	render.Output(os.Stdout, &jobsOutputer{}, out.Jobs, jobsListArgs.Output)
	return nil
}

type jobsOutputer struct{}

func (o *jobsOutputer) Standard(w io.Writer, generic interface{}) error {
	outs := generic.([]modzysdkmodel.JobDetails)

	tabbed := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	fmt.Fprintf(tabbed, "ID\tStatus\tSubmitted By\tTeam\tSubmitted At\tModel\n")
	fmt.Fprintf(tabbed, "--\t------\t------------\t----\t------------\t-----\n")

	for _, out := range outs {
		fmt.Fprintf(tabbed, "%s\t%s\t%s\t%s\t%s\t%s\n",
			out.JobIdentifier, out.Status, out.SubmittedBy, out.Team.Identifier,
			out.SubmittedAt, out.Model.Name,
		)
	}
	if err := tabbed.Flush(); err != nil {
		return err
	}

	return nil
}
