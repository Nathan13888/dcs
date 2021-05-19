package cmd

import (
	"dcs/config"
	"dcs/scraper"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dcs",
	Short: "A Golang scraper for Dramacool/WatchAsian",
	Long: `A featureful scraper for Dramacool/WatchAsian.

Written in Go with Colly, Cobra and several other libraries.

Could be configured and also features a daemon to periodically check for new drama episodes.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	url, err := rootCmd.PersistentFlags().GetString("url")
	if err != nil {
		fmt.Println(err)
	}
	newUrl := strings.Trim(url, " /")
	if len(newUrl) > 0 {
		fmt.Printf("New URL is not set to `%s`\n\n", newUrl)
		scraper.URL = newUrl
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringP("url", "u", "", "specify DC url")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dcs.json)")
	rootCmd.PersistentFlags().BoolP("remote", "r", false, "Download from remote server")
	rootCmd.PersistentFlags().StringP("remote-host", "a", "", "Specify host address of remote DCS")
	rootCmd.PersistentFlags().StringP("remote-port", "p", "", "Specify port of remote DCS")
	viper.BindPFlag("DaemonHost", rootCmd.PersistentFlags().Lookup("remote-host"))
	viper.BindPFlag("DaemonPort", rootCmd.PersistentFlags().Lookup("remote-port"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		config.Configure()
	}

	viper.AutomaticEnv() // read in environment variables that match

	config.ConfigRecents()
}
