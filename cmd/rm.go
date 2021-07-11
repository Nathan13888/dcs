package cmd

import (
	"dcs/server"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

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
		// purge, err := cmd.Flags().GetBool("purge")
		// if err != nil {
		// 	panic(err)
		// }
		for _, x := range args {
			if strings.EqualFold(x, "purge") {
				res, err := Request("GET", "api/purgejobs")
				if err != nil {
					panic(err)
				}
				defer res.Body.Close()
				var jobs server.JobsResponse
				decoder := json.NewDecoder(res.Body)
				decoder.DisallowUnknownFields()
				err = decoder.Decode(&jobs)
				if err != nil {
					panic(err)
				}
				fmt.Println("Purged the following jobs")
				return
			}

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

	// rmCmd.Flags().BoolP("purge", "P", false, "purge all broken or unreliable jobs")
}
