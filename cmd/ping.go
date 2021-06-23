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
	Aliases: []string{"p"},
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
}
