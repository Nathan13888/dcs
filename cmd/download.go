package cmd

import (
	"dcs/scraper"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: sanitize arguments
		var link string
		if scraper.IsLink(args[0]) {
			link = args[0]
		} else {
			res := scraper.FirstSearch(scraper.JoinArgs(args[:len(args)-1]))
			if res != "" {
				// TODO: sanitize value of episode
				num, err := strconv.Atoi(args[len(args)-1])
				if err != nil {
					panic(err)
				}
				episodes := scraper.GetEpisodesByLink(res)
				var url string
				for _, e := range episodes {
					if e.Number == num {
						url = e.Link
					}
				}
				if len(episodes) >= num || url == "" {
					link = scraper.URL + url
				} else {
					fmt.Printf("Episode %d was not available", num)
				}
			} else {
				fmt.Println("There has been a problem using your specified query")
				return
			}
		}
		ajax := scraper.GetAjax(link)
		if ajax.Found {
			fmt.Printf("Attemping to download an episode from '%s'\n\n", link)
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
			fmt.Println("FAILED to find episode...")
		}
	},
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
