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
		fmt.Println("Version:\t", config.BuildVersion)
		fmt.Println("Built by:\t", config.BuildUser)
		fmt.Println("Build Time:\t", config.BuildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
