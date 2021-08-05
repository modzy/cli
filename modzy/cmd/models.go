package cmd

import (
	"github.com/spf13/cobra"
)

var modelsArgs struct {
}

func init() {
	rootCmd.AddCommand(modelsCmd)
}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Work with models",
	Long:  ``,
}
