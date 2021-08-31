package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/modzy/cli/internal/render"
	"github.com/modzy/cli/internal/status"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

var projectStatusArgs struct {
	Output string
}

func init() {
	projectStatusCmd.Flags().StringVarP(&projectStatusArgs.Output, "output", "o", "", "")

	projectsCmd.AddCommand(projectStatusCmd)
}

var projectStatusCmd = &cobra.Command{
	Use:          "status [projectID]",
	Short:        "Returns project level dashboard information",
	Args:         cobra.ExactArgs(1),
	RunE:         projectStatusRun,
	SilenceUsage: true,
}

func projectStatusRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	fetcher := status.NewFetcher(client)

	stats := fetcher.FetchProject(ctx, args[0])

	render.Output(os.Stdout, &projectStatusOutputer{}, &stats, projectStatusArgs.Output)
	return nil
}

type projectStatusOutputer struct{}

func (o *projectStatusOutputer) Standard(w io.Writer, generic interface{}) error {
	stats := generic.(*status.ProjectModel)

	// errors
	if len(stats.Errors) > 0 {
		if _, err := fmt.Fprintf(w, "Errors:\n"); err != nil {
			// for testing coverage
			return err
		}
		for _, err := range stats.Errors {
			fmt.Fprintf(w, "- %s", err)
		}
	}

	// processed
	fmt.Fprintf(w, "Data processed: %s\n",
		humanize.Bytes(uint64(stats.DataProcessed.Summary.RecentBytes)),
	)

	// predictions-made
	fmt.Fprintf(w, "Predictions made: %d\n",
		stats.Predictions.Summary.RecentPredictions,
	)

	// active-models
	fmt.Fprintf(w, "Active Models:\n")
	if len(stats.ActiveModels) == 0 {
		fmt.Fprintf(w, "  No active models\n")
	} else {
		tabbed := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabbed, "Rank\tName\tVersion\tElapsed\n")
		fmt.Fprintf(tabbed, "----\t----\t-------\t-------\n")
		for _, a := range stats.ActiveModels {
			fmt.Fprintf(tabbed, "%d\t%s\t%s\t%s\n", a.Ranking, a.Name, a.Version, time.Duration(a.ElapsedTime)*time.Millisecond)
		}
		tabbed.Flush()
		fmt.Fprintf(w, "\n")
	}

	return nil
}
