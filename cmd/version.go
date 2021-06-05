package cmd

import (
	"dcs/config"
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Shows versioning info about DCS",
	Long:    `Shows versioning info about DCS`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version:\t %s (%s/%s)\n", config.BuildVersion, config.BuildGOOS, config.BuildARCH)
		fmt.Println("Built by:\t", config.BuildUser)
		fmt.Println("Build Time:\t", config.BuildTime)
		fmt.Printf("Running on:\t %s/%s\n", config.GOOS, config.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
