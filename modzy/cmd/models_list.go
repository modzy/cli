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
	pagingArgs
	Output string
}

func init() {
	modelsListCmd.Flags().StringArrayVarP(&modelsListArgs.Filter, "filter", "", []string{}, "")
	modelsListCmd.Flags().IntVarP(&modelsListArgs.Take, "take", "", 10, "")
	modelsListCmd.Flags().IntVarP(&modelsListArgs.Page, "page", "", 1, "")
	modelsListCmd.Flags().StringVarP(&modelsListArgs.Output, "output", "o", "", "")

	modelsCmd.AddCommand(modelsListCmd)
}

var modelsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List models",
	Long:  ``,
	RunE:  modelsListRun,
}

func modelsListRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	paging, err := modelsListArgs.GetPagingInput()
	if err != nil {
		return err
	}

	input := &modzysdk.ListModelsInput{
		Paging: paging,
	}
	out, err := client.Models().ListModels(ctx, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	summaries := []modelSummaryWithMore{}
	for _, modelSummary := range out.Models {
		detail, err := client.Models().GetModelDetails(ctx, &modzysdk.GetModelDetailsInput{
			ModelID: modelSummary.ID,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
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
