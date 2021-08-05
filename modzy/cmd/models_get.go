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

var modelsGetArgs struct {
	Output string
}

func init() {
	modelsGetCmd.Flags().StringVarP(&modelsGetArgs.Output, "output", "o", "", "TODO: good output description")

	modelsCmd.AddCommand(modelsGetCmd)
}

var modelsGetCmd = &cobra.Command{
	Use:   "get [modelID]",
	Short: "Get detailed information about a model",
	Long:  `This will get effective authentication information that this command is using and will display the non-sensitive portion for the purpose of troubleshooting authentication.`,
	Args:  cobra.ExactArgs(1),
	RunE:  modelsGetRun,
}

func modelsGetRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	modelID := args[0]
	out, err := client.Models().GetModelDetails(ctx, &modzysdk.GetModelDetailsInput{
		ModelID: modelID,
	})
	if err != nil {
		return err
	}

	render.Output(os.Stdout, &ModelOutputer{}, out, modelsGetArgs.Output)
	return nil
}

type ModelOutputer struct{}

func (o *ModelOutputer) RowHeader(w io.Writer) error {
	fmt.Fprintf(w, "ID\tLatest Version\tAuthor\n")
	fmt.Fprintf(w, "--\t--------------\t------\n")
	return nil
}

func (o *ModelOutputer) RowData(w io.Writer, generic interface{}) error {
	out := generic.(*modzysdk.GetModelDetailsOutput)
	fmt.Fprintf(w, "%s\t%s\t%s\n", out.Details.ModelID, out.Details.LatestVersion, out.Details.Author)
	return nil
}

func (o *ModelOutputer) Standard(w io.Writer, generic interface{}) error {
	out := generic.(*modzysdk.GetModelDetailsOutput)

	tabbed := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(tabbed, "ID: \t%s\n", out.Details.ModelID)
	fmt.Fprintf(tabbed, "Name: \t%s\n", out.Details.Name)
	fmt.Fprintf(tabbed, "Author: \t%s\n", out.Details.Author)
	fmt.Fprintf(tabbed, "Versions: \t%s\n", strings.Join(out.Details.Versions, ", "))
	fmt.Fprintf(tabbed, "Description: \t%s\n", out.Details.Description)

	if err := tabbed.Flush(); err != nil {
		return err
	}

	return nil
}
