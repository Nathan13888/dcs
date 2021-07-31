package cmd

import (
	"dcs/config"
	"dcs/prompt"
	"dcs/scraper"
	"fmt"
	"os"

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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Test if err != nil, then ask user whether to proceed with error
func testError(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	exit := !prompt.Confirm("Would you like to continue?")
	if exit {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&scraper.URL, "url", "u", "https://watchasian.cc", "specify DC url")
	rootCmd.PersistentFlags().BoolVarP(&scraper.ASIANLOAD, "asianload", "A", false, "Use Asianload instead")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dcs.json)")
	rootCmd.PersistentFlags().BoolP("remote", "r", false, "Download from remote server")
	rootCmd.PersistentFlags().StringP("remote-host", "a", "", "Specify host address of remote DCS")
	rootCmd.PersistentFlags().StringP("remote-port", "p", "", "Specify port of remote DCS")
	rootCmd.PersistentFlags().BoolVarP(&scraper.NOPROXY, "no-proxy", "N", false, "disable all proxies")
	viper.BindPFlag("DaemonHost", rootCmd.PersistentFlags().Lookup("remote-host"))
	viper.BindPFlag("DaemonPort", rootCmd.PersistentFlags().Lookup("remote-port"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// TODO: verify new url is functional
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

	// fmt.Println(scraper.ASIANLOAD)
	// fmt.Println(scraper.URL)
	if scraper.ASIANLOAD {
		scraper.URL = "https://asianload.io"
	}
}
