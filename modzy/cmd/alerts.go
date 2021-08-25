package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/modzy/cli/internal/render"
	modzy "github.com/modzy/sdk-go"

	"github.com/spf13/cobra"
)

var alertsGetArgs struct {
	Type   string
	Output string
}

func init() {
	alertsGetCmd.Flags().StringVarP(&alertsGetArgs.Output, "output", "o", "", "")
	rootCmd.AddCommand(alertsGetCmd)
}

var alertsGetCmd = &cobra.Command{
	Use:          "alerts [type]",
	Short:        "List all alerts",
	Long:         ``,
	RunE:         alertsGetRun,
	SilenceUsage: true,
}

func alertsGetRun(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	client := getClient()

	out, err := client.Dashboard().GetAlerts(ctx, &modzy.GetAlertsInput{})
	if err != nil {
		return err
	}

	details := []alertAndDetail{}

	for _, asum := range out.Alerts {
		atype := string(asum.Type)
		alert, err := client.Dashboard().GetAlertDetails(ctx, &modzy.GetAlertDetailsInput{Type: asum.Type})
		if err != nil {
			return err
		}
		details = append(details, alertAndDetail{
			Type:     atype,
			Count:    asum.Count,
			Entities: alert.Entities,
		})
	}

	render.Output(os.Stdout, &alertsGetOutputer{}, details, alertsGetArgs.Output)
	return nil
}

type alertsGetOutputer struct{}

func (o *alertsGetOutputer) Standard(w io.Writer, generic interface{}) error {
	alerts := generic.([]alertAndDetail)

	if len(alerts) == 0 {
		fmt.Printf("No alerts\n")
	} else {
		for _, aAndD := range alerts {
			fmt.Printf("%s: %d\n", aAndD.Type, aAndD.Count)
			for _, ent := range aAndD.Entities {
				fmt.Printf("- %s\n", ent)
			}
		}
	}

	return nil
}

type alertAndDetail struct {
	Type     string
	Count    int
	Entities []string
}
