package cmd

import (
	"dcs/scraper"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Gives you information about a drama.",
	Long: `Know about the episodes and other info related to anything you want on DramaCool/WatchAsian.

	USAGE: dcs info <link here or name of drama>`,
	Run: func(cmd *cobra.Command, args []string) {
		var link string
		if !strings.Contains(scraper.JoinArgs(args), "/") {
			res := scraper.Search(args[0])
			if len(res) > 0 {
				link = res[0].FullURL
			}
		} else {
			link = args[0]
		}
		res := scraper.GetEpisodesByLink(link)
		fmt.Printf("Displaying info for '%s'\n", link)
		// TODO: show info about description
		fmt.Printf("\n%d Episodes Available\n\n", len(res))
		for _, e := range res {
			fmt.Printf("Episode %d was available on %s --> %s\n",
				e.Number, e.Date, scraper.URL+e.Link)
		}
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
