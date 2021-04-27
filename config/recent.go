package config

import (
	"dcs/scraper"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var configHome string
var recentsConfigPath string
var timeFormat = "Jan-02-06_15:04:05"

var recentsConfig = viper.New()
var recentDownloads map[string][]string = make(map[string][]string)
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
		recentDownloads = recentsConfig.GetStringMapStringSlice("downloads")
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

func AddRecentDownload(info *scraper.DramaInfo) {
	recentDownloads[strings.ToLower(info.Name)] = []string{
		time.Now().Format(timeFormat),
		info.Name,
		info.SubURL,
	}
	SyncRecents()
}

func AddRecentSearch(s string) {
	recentSearches = append(recentSearches, s)
	SyncRecents()
}

func GetRecentDownloads() []scraper.DramaInfo {
	var dramas []scraper.DramaInfo
	for key, props := range recentDownloads {
		if !(len(props) >= 2) {
			panic(fmt.Errorf("invalid properties for key `%s`: %s", key, props))
		}
		name := props[1]
		subURL := props[2]

		dramas = append(dramas, scraper.DramaInfo{
			FullURL: scraper.URL + subURL,
			SubURL:  subURL,
			Domain:  scraper.URL,
			Name:    name,
		})
	}

	return dramas
}

func GetRecentSearches() []string {
	return recentSearches
}
