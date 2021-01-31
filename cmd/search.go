package cmd

import (
	"fmt"
	"math"
	"strings"

	"dcs/scraper"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "USAGE: dcs search <search phrase>",
	Long:  `Search for whatever you want...`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Save results as context/search

		qry := strings.ReplaceAll(strings.Join(args, " "), "\"", "")
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
