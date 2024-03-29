package cmd

import (
	"fmt"
	"math"

	"dcs/scraper"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a drama",
	Long: `Search for anything on the website.

	USAGE: dcs search <search phrase>`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Save results as context/search

		qry := scraper.JoinArgs(args)
		fmt.Printf("Searching for '%s'...\n\n", qry)
		res := scraper.Search(qry)
		if len(res) == 0 {
			fmt.Println("Sorry, no results found...")
		} else {
			fmt.Printf("Here are the first 10 results:\n\n")
			for i := 0; i < (int)(math.Min(10, (float64)(len(res)))); i++ {
				var option scraper.DramaInfo = res[i]
				fmt.Printf("%d) %s --> %s\n\n", i+1, option.Name, option.FullURL)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
