package cmd

import (
	"github.com/spf13/cobra"
)

// var jobsArgs struct {
// }

func init() {
	rootCmd.AddCommand(jobsCmd)
}

var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Work with jobs",
	Long:  ``,
}
