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

var statusArgs struct {
	Output string
}

func init() {
	statusCmd.Flags().StringVarP(&statusArgs.Output, "output", "o", "", "")

	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:          "status",
	Short:        "Returns top level dashboard information for your account",
	RunE:         statusRun,
	SilenceUsage: true,
}

func statusRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	fetcher := status.NewFetcher(client)

	top := fetcher.FetchTop(ctx)

	render.Output(os.Stdout, &statusOutputer{}, &top, statusArgs.Output)
	return nil
}

type statusOutputer struct{}

func (o *statusOutputer) Standard(w io.Writer, generic interface{}) error {
	top := generic.(*status.TopModel)

	// errors
	if len(top.Errors) > 0 {
		if _, err := fmt.Fprintf(w, "Errors:\n"); err != nil {
			// for testing coverage
			return err
		}
		for _, err := range top.Errors {
			fmt.Fprintf(w, "- %s", err)
		}
	}

	// alerts
	fmt.Fprintf(w, "Alerts:\n")
	if len(top.Alerts) == 0 {
		fmt.Fprintf(w, "  None\n")
	}
	for _, a := range top.Alerts {
		fmt.Fprintf(w, "- %s (%d)\n", a.Type, a.Count)
	}

	// processed
	fmt.Fprintf(w, "Data processed: %s\n",
		humanize.Bytes(uint64(top.DataProcessed.Summary.RecentBytes)),
	)

	// predictions-made
	fmt.Fprintf(w, "Predictions made: %d\n",
		top.Predictions.Summary.RecentPredictions,
	)

	// accounting-users
	fmt.Fprintf(w, "\nTotal Users: %d\n", len(top.Users))

	// active-users
	fmt.Fprintf(w, "Active Users:\n")
	if len(top.ActiveUsers) == 0 {
		fmt.Fprintf(w, "  No active users\n")
	} else {
		tabbed := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabbed, "Rank\tName\tElapsed\n")
		fmt.Fprintf(tabbed, "----\t----\t-------\n")
		for _, a := range top.ActiveUsers {
			fmt.Fprintf(tabbed, "%d\t%s\t%s\n", a.Ranking, fmt.Sprintf("%s %s", a.FirstName, a.LastName), time.Duration(a.ElapsedTime)*time.Millisecond)
		}
		tabbed.Flush()
		fmt.Fprintf(w, "\n")
	}

	// active-models
	fmt.Fprintf(w, "Active Models:\n")
	if len(top.ActiveModels) == 0 {
		fmt.Fprintf(w, "  No active models\n")
	} else {
		tabbed := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
		fmt.Fprintf(tabbed, "Rank\tName\tVersion\tElapsed\n")
		fmt.Fprintf(tabbed, "----\t----\t-------\t-------\n")
		for _, a := range top.ActiveModels {
			fmt.Fprintf(tabbed, "%d\t%s\t%s\t%s\n", a.Ranking, a.Name, a.Version, time.Duration(a.ElapsedTime)*time.Millisecond)
		}
		tabbed.Flush()
		fmt.Fprintf(w, "\n")
	}

	// latest models
	fmt.Fprintf(w, "Latest Models : %d\n", len(top.LatestModels))

	// cpu-overall-usage
	fmt.Fprintf(w, "\nCPU usage:\n")
	fmt.Fprintf(w, "- Average: %0.0f\n", top.CPUOverallAverage)
	fmt.Fprintf(w, "- Last: %0.0f\n", top.CPUOverallLast)

	// engines
	fmt.Fprintf(w, "Engines: %d of %s\n", top.EnginesProcessing, top.License.ProcessingEngines)

	return nil
}
