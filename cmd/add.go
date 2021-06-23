package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a drama to periodically download",
	Long: `Add a drama for the DCS daemon to periodically check for updates and download.

	USAGE: service add <link to drama>
	USAGE: service add <name of drama>`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
}

func init() {
	serviceCmd.AddCommand(addCmd)
}
