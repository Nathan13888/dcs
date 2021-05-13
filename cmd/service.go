package cmd

import (
	"dcs/config"
	"dcs/scraper"
	"fmt"

	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Control the DCS daemon",
	Long:  `A daemon to periodically download new episodes of a drama.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("service called")
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GetRemoteURL(rpath string) string {
	ip, port := config.DaemonURL()
	url := scraper.JoinURL(fmt.Sprintf("http://%s:%d", ip, port), rpath)
	return url
}
