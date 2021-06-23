package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Change a config setting",
	Long:  `It's complicated...`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		prop := args[0]
		fmt.Printf("Config property `%s` is set to `%s`\n", prop, viper.GetString(prop))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
