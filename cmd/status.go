package cmd

import (
	"dcs/config"
	"dcs/scraper"
	"fmt"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status info about the DCS daemon",
	Long: `Info about the DCS daemon...

	USAGE: service status`,
	Run: func(cmd *cobra.Command, args []string) {
		host, port := config.DaemonURL()
		fmt.Printf("Displaying information about `%s:%d`\n\n", host, port)

		// TODO: change to pinging API
		online := scraper.Ping(host, port)
		// TODO: display extra stats about online server
		if online {
			fmt.Println("Status: Online")
		} else {
			fmt.Println("Status: Offline")
		}
	},
}

func init() {
	serviceCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
