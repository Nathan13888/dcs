package cmd

import (
	"dcs/config"
	"dcs/scraper"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Control the DCS daemon",
	Long:  `A daemon to periodically download new episodes of a drama.`,
	Aliases: []string{
		"s",
		"r",
		"rem",
		"remote",
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}

func Request(method string, rpath string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, GetRemoteURL(rpath), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)

	return res, err
}

func GetRemoteURL(rpath string) string {
	ip, port := config.DaemonURL()
	url := scraper.JoinURL(fmt.Sprintf("http://%s:%d", ip, port), rpath)
	return url
}
