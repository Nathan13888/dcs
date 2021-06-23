package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a drama from being downloaded periodically",
	Long: `Remove a drama for the DCS daemon to periodically check for updates and download.

	USAGE: service rm <id of drama in list>`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, x := range args {
			id := url.PathEscape(x)
			fmt.Printf("Removing job %s\n", id)
			res, err := Request("DELETE", "api/remove/"+id)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Response: %s\n", res.Status)
		}
	},
}

func init() {
	serviceCmd.AddCommand(rmCmd)
}
