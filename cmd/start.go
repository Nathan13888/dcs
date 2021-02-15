package cmd

import (
	"dcs/daemon"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the DCS daemon server",
	Long: `Start the DCS daemon, a HTTP REST api.

	Usage: service start`,
	Run: func(cmd *cobra.Command, args []string) {
		daemon.Start()
	},
}

func init() {
	serviceCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
