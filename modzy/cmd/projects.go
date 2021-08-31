package cmd

import (
	"github.com/spf13/cobra"
)

// var projectsArgs struct {
// }

func init() {
	rootCmd.AddCommand(projectsCmd)
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Work with projects",
	Long:  ``,
}
