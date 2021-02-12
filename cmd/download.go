package cmd

import (
	"dcs/scraper"
	"fmt"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download an episode or episodes of a drama",
	Long: `Download anything from DCS that you want.

	USAGE: download <link to episode>
	USAGE: download <name of drama> <episode range>`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: sanitize arguments
		if scraper.IsLink(args[0]) {
			download(args[0])
		} else {
			res := scraper.FirstSearch(scraper.JoinArgs(args[:len(args)-1]))
			if res != "" {
				// TODO: sanitize value of episode
				episodes := scraper.GetEpisodesByLink(res)
				episodeRange := scraper.GetRange(args[len(args)-1])
				fmt.Printf("Attemping to download these episodes: %v\n\n", episodeRange)
				for _, num := range episodeRange {
					fmt.Printf("Looking for episode %d...\n", num)
					var url string
					for _, e := range episodes {
						if e.Number == num {
							url = e.Link
						}
					}
					if len(episodes) >= num || url == "" {
						download(scraper.URL + url)
					} else {
						fmt.Printf("Episode %d was not available", num)
					}
				}
			} else {
				fmt.Println("There has been a problem using your specified query")
				return
			}
		}
	},
}

func download(link string) {
	ajax := scraper.GetAjax(link)
	if ajax.Found {
		fmt.Printf("Attemping to download from '%s'\n\n", link)
		fmt.Printf("Found AJAX endpoint '%s'\n\n", ajax.Ajax)
		link := scraper.ScrapeEpisode(ajax)
		fmt.Printf("Found '%s'\n\n", link)
		// TODO: prompt confirm download
		fmt.Println("Downloading...")
		err := scraper.Download(scraper.DownloadInfo{
			Link: link,
			Name: ajax.Name,
			Num:  ajax.Num,
		})
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Print("FAILED to find episode...\n\n")
	}
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
