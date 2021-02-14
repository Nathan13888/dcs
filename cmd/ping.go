package cmd

import (
	"dcs/scraper"
	"fmt"

	"github.com/spf13/cobra"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping Dramacool/Watchasian to see if their website is up",
	Long: `Pings Dramacool/Watchasian to see if their website is up.

	USAGE: ping`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Pinging %s...\n", scraper.URL)
		res := scraper.PingDC()
		if res {
			fmt.Println("The ping was *successful*")
		} else {
			fmt.Println("The ping was *unsuccessful*")
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
