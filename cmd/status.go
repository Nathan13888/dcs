package cmd

import (
	"dcs/config"
	"dcs/server"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get status info about the DCS daemon",
	Long: `Info about the DCS daemon...

	USAGE: service status`,
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		host, port := config.DaemonURL()
		fmt.Printf("Displaying information about `%s:%d`\n\n", host, port)

		res, err := Request("GET", "status")
		if err != nil {
			// fmt.Println(err)
			goto DeclareOffline
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			fmt.Printf("*** ONLINE ***")
			var obj server.StatusResponse

			decoder := json.NewDecoder(res.Body)
			decoder.DisallowUnknownFields()
			err = decoder.Decode(&obj)
			if err != nil {
				panic(err)
			}

			fmt.Printf("\nUptime:          \t%s\n", obj.Uptime.Round(time.Second).String())
			fmt.Printf("Dramas:            \t%d\n", obj.DownloadedDramas)
			fmt.Printf("Episodes:          \t%d\n", obj.DownloadedEpisodes)
			fmt.Printf("Library Size:      \t%.3f GBs\n", float64(obj.LibrarySize)/math.Pow(1024, 3))
			fmt.Printf("Processed Requests:\t%d\n", obj.ProcessedRequests)
			fmt.Printf("DCS Version:       \t%s\n", obj.Version)
			fmt.Printf("Build Info:        \t%s\n", obj.BuildInfo)

			return
		}
	DeclareOffline:
		fmt.Println("*** OFFLINE ***")
	},
}

func init() {
	serviceCmd.AddCommand(statusCmd)
}
