package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/modzy/cli/modzy/render"
	modzysdk "github.com/modzy/sdk-go"

	"github.com/spf13/cobra"
)

var modelsListArgs struct {
	Filter []string
	Take   int
	Page   int
	Output string
}

func init() {
	modelsListCmd.Flags().StringArrayVarP(&modelsListArgs.Filter, "filter", "", []string{}, "TODO: good description")
	modelsListCmd.Flags().IntVarP(&modelsListArgs.Take, "take", "", 10, "TODO: good description")
	modelsListCmd.Flags().IntVarP(&modelsListArgs.Page, "page", "", 1, "TODO: good description")
	modelsListCmd.Flags().StringVarP(&modelsListArgs.Output, "output", "o", "", "TODO: good description")

	modelsCmd.AddCommand(modelsListCmd)
}

var modelsListCmd = &cobra.Command{
	Use:   "list [modelID]",
	Short: "List models information about a model",
	Long:  ``,
	RunE:  modelsListRun,
}

func modelsListRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	input := (&modzysdk.ListModelsInput{}).
		WithPaging(modelsListArgs.Take, modelsListArgs.Page)

	for _, filter := range modelsListArgs.Filter {
		filterSplit := strings.Split(filter, "=")
		if len(filterSplit) != 2 {
			return fmt.Errorf("Filter is not correctly formatted (field=value): '%s'", filter)
		}
		input.Paging = input.Paging.WithFilterAnd(filterSplit[0], filterSplit[1])
	}

	out, err := client.Models().ListModels(ctx, input)
	if err != nil {
		return err
	}

	summaries := []modelSummaryWithMore{}
	for _, modelSummary := range out.Models {
		detail, err := client.Models().GetModelDetails(ctx, &modzysdk.GetModelDetailsInput{
			ModelID: modelSummary.ID,
		})
		if err != nil {
			return err
		}
		summaries = append(summaries, modelSummaryWithMore{
			ID:       modelSummary.ID,
			Name:     detail.Details.Name,
			Author:   detail.Details.Author,
			Versions: detail.Details.Versions,
		})
	}
	render.Output(os.Stdout, &ModelsOutputer{}, summaries, modelsListArgs.Output)
	return nil
}

type ModelsOutputer struct{}

func (o *ModelsOutputer) Standard(w io.Writer, generic interface{}) error {
	outs := generic.([]modelSummaryWithMore)

	tabbed := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	fmt.Fprintf(tabbed, "ID\tName\tAuthor\tVersions\n")
	fmt.Fprintf(tabbed, "--\t----\t------\t--------\n")

	for _, out := range outs {
		fmt.Fprintf(tabbed, "%s\t%s\t%s\t%s\n",
			out.ID, out.Name, out.Author,
			strings.Join(out.Versions, ", "),
		)
	}
	if err := tabbed.Flush(); err != nil {
		return err
	}

	return nil
}

type modelSummaryWithMore struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Author        string   `json:"author"`
	LatestVersion string   `json:"latestVersion"`
	Versions      []string `json:"versions"`
}
