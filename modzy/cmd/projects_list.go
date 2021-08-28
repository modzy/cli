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

var projectsListArgs struct {
	pagingArgs
	Output string
}

func init() {
	projectsListCmd.Flags().StringArrayVarP(&projectsListArgs.Filter, "filter", "", []string{}, "")
	projectsListCmd.Flags().IntVarP(&projectsListArgs.Take, "take", "", 10, "")
	projectsListCmd.Flags().IntVarP(&projectsListArgs.Page, "page", "", 1, "")
	projectsListCmd.Flags().StringVarP(&projectsListArgs.Output, "output", "o", "", "")

	projectsCmd.AddCommand(projectsListCmd)
}

var projectsListCmd = &cobra.Command{
	Use:          "list",
	Short:        "List projects",
	Long:         ``,
	RunE:         projectsListRun,
	SilenceUsage: true,
}

func projectsListRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	paging, err := projectsListArgs.GetPagingInput()
	if err != nil {
		return err
	}

	input := &modzysdk.ListProjectsInput{
		Paging: paging,
	}
	out, err := client.Accounting().ListProjects(ctx, input)
	if err != nil {
		return err
	}

	render.Output(os.Stdout, &projectsListOutputer{}, out.Projects, projectsListArgs.Output)
	return nil
}

type projectsListOutputer struct{}

func (o *projectsListOutputer) Standard(w io.Writer, generic interface{}) error {
	outs := generic.([]modzysdkmodel.AccountingProject)

	tabbed := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)

	fmt.Fprintf(tabbed, "ID\tName\tOwner\tKey Prefix\tDate Created\n")
	fmt.Fprintf(tabbed, "--\t----\t-----\t----------\t------------\n")

	for _, out := range outs {
		keyPrefix := ""
		if len(out.AccessKeys) > 0 {
			keyPrefix = out.AccessKeys[0].Prefix
		}
		fmt.Fprintf(tabbed, "%s\t%s\t%s\t%s\t%s\n",
			out.Identifier, out.Name, fmt.Sprintf("%s %s", out.User.FirstName, out.User.LastName), keyPrefix, out.CreatedAt,
		)
	}
	if err := tabbed.Flush(); err != nil {
		return err
	}

	return nil
}
