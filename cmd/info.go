package cmd

import (
	"dcs/scraper"
	"fmt"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Gives you information about a drama",
	Long: `Know about the episodes and other info related to anything you want on DramaCool/WatchAsian.

	USAGE: dcs info <link to drama>
	USAGE: dcs info <name of drama>`,
	Run: func(cmd *cobra.Command, args []string) {
		var link string
		if scraper.IsLink(args[0]) {
			link = args[0]
		} else {
			link = scraper.FirstSearch(scraper.JoinArgs(args))
		}
		res := scraper.GetEpisodesByLink(link)
		fmt.Printf("Displaying info for '%s'\n", link)
		// TODO: show info about description
		fmt.Printf("\n%d Episodes Available\n\n", len(res))
		DisplayEpisodesInfo(res)
	},
}

// DisplayEpisodesInfo - Lists all all the information about a list of episodes
func DisplayEpisodesInfo(episodes []scraper.EpisodeInfo) {
	for _, e := range episodes {
		fmt.Printf("Episode %v was available on %s --> %s\n",
			e.Number, e.Date, scraper.URL+e.Link)
	}
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
