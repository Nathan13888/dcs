package cmd

import (
	"dcs/config"
	"dcs/downloader"
	"dcs/prompt"
	"dcs/scraper"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download an episode or episodes of a drama",
	Long: `Download anything from DCS that you want.

	USAGE: download  -->  (for interactive prompt)
	USAGE: download <link to episode>
	USAGE: download <name of drama> <episode range>`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: sanitize arguments
		if len(args) == 1 && scraper.IsLink(args[0]) {
			download(args[0])
		} else {
			var link string
			var episodeRange []int
			if len(args) == 0 {
				var res string
				var err error
				var drama scraper.DramaInfo
				res, err = prompt.String("Search")
				if err != nil {
					panic(err)
				}
				// recent := config.GetRecentDownloads()
				queries := scraper.Search(res)
				if len(queries) == 0 {
					// TODO: don't PANIC
					panic("no queries were found")
				} else {
					resInfo, err := prompt.Drama(queries)
					if err != nil {
						panic(err)
					}
					drama = *resInfo
					//TODO: more rigorous checking
					link = drama.FullURL
				}

				config.AddRecentDownload(drama.SubURL)

				episodes := scraper.GetEpisodes(drama)
				DisplayEpisodesInfo(episodes)

				res, err = prompt.String("Episode Range")
				if err != nil {
					panic(err)
				}
				episodeRange = scraper.GetRange(strings.TrimSpace(res))
			} else {
				link = scraper.FirstSearch(scraper.JoinArgs(args[:len(args)-1]))
				episodeRange = scraper.GetRange(args[len(args)-1])
			}
			if link != "" {
				// TODO: sanitize value of episode
				episodes := scraper.GetEpisodesByLink(link)
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
		link := scraper.ScrapeAjax(ajax)
		fmt.Printf("Found '%s'\n\n", link)
		// TODO: prompt confirm download
		fmt.Println("Downloading...")
		err := downloader.Get(downloader.DownloadInfo{
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
