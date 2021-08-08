package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/modzy/cli/modzy/render"
	modzysdk "github.com/modzy/sdk-go"
	modzysdkmodel "github.com/modzy/sdk-go/model"

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
	Long:  ``,
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

	render.Output(os.Stdout, &ModelOutputer{}, out.Details, modelsGetArgs.Output)
	return nil
}

type ModelOutputer struct{}

func (o *ModelOutputer) Standard(w io.Writer, generic interface{}) error {
	out := generic.(modzysdkmodel.ModelDetails)

	tabbed := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(tabbed, "ID: \t%s\n", out.ModelID)
	fmt.Fprintf(tabbed, "Name: \t%s\n", out.Name)
	fmt.Fprintf(tabbed, "Author: \t%s\n", out.Author)
	fmt.Fprintf(tabbed, "Versions: \t%s\n", strings.Join(out.Versions, ", "))
	fmt.Fprintf(tabbed, "Description: \t%s\n", out.Description)

	if err := tabbed.Flush(); err != nil {
		return err
	}

	return nil
}
