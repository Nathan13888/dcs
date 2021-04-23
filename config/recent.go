package config

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var configHome string
var recentsConfigPath string

var recentsConfig = viper.New()
var recentDownloads []string = make([]string, 0)
var recentSearches []string = make([]string, 0)

func ConfigRecents() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".dcs" (without extension).
	configHome = path.Join(home, ".dcs")
	recentsConfigPath = path.Join(configHome, "recent.json")
	recentsConfig.SetConfigName("recent")
	recentsConfig.SetConfigType("json")
	recentsConfig.AddConfigPath(configHome)

	// If a config file is found, read it in.
	if err = recentsConfig.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", recentsConfig.ConfigFileUsed())
		recentDownloads = recentsConfig.GetStringSlice("downloads")
		recentSearches = recentsConfig.GetStringSlice("searches")
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		os.MkdirAll(configHome, 0755)
		recentsConfig.SafeWriteConfigAs(recentsConfigPath)
	} else {
		panic(err)
	}
	recentsConfig.SetDefault("downloads", []string{})
	recentsConfig.SetDefault("searches", []string{})
}

func SyncRecents() {
	recentsConfig.Set("downloads", recentDownloads)
	recentsConfig.Set("searches", recentSearches)
	recentsConfig.WriteConfigAs(recentsConfigPath)
}

func AddRecentDownload(subURL string) {
	// fmt.Println(subURL)
	recentDownloads = append(recentDownloads, subURL)
	SyncRecents()
}

func AddRecentSearch(s string) {
	recentSearches = append(recentSearches, s)
	SyncRecents()
}

func GetRecentDownloads() []string {
	return recentDownloads
}

func GetRecentSearches() []string {
	return recentSearches
}
