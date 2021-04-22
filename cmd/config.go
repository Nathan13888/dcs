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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			prop := args[0]
			fmt.Printf("Config property `%s` is set to `%s`\n", prop, viper.GetString(prop))
		} else {
			fmt.Println("This command currently only displays config values.")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
