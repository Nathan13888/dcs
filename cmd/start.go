package cmd

import (
	"dcs/server"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the DCS daemon server",
	Long: `Start the DCS daemon, a HTTP REST api.

	Usage: service start`,
	Run: func(cmd *cobra.Command, args []string) {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			panic(err)
		}
		server.Start(debug)
	},
}

func init() {
	serviceCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("debug", "d", false, "Debug mode")
}
