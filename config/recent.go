package config

import (
	"dcs/scraper"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var recentsConfigPath string
var timeFormat = "Jan-02-06_15:04:05"

var recentsConfig = viper.New()
var recentDownloads map[string][]string = make(map[string][]string)

// var recentSearches []string = make([]string, 0)

func ConfigRecents() {
	// Search config in home directory with name ".dcs" (without extension).
	recentsConfigPath = path.Join(GetConfigHome(), "recent.json")
	recentsConfig.SetConfigName("recent")
	recentsConfig.SetConfigType("json")
	recentsConfig.AddConfigPath(GetConfigHome())

	// If a config file is found, read it in.
	if err := recentsConfig.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", recentsConfig.ConfigFileUsed())
		recentDownloads = recentsConfig.GetStringMapStringSlice("downloads")
		// recentSearches = recentsConfig.GetStringSlice("searches")
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		os.MkdirAll(GetConfigHome(), 0755)
		recentsConfig.SafeWriteConfigAs(recentsConfigPath)
	} else {
		panic(err)
	}
	recentsConfig.SetDefault("downloads", []string{})
	// recentsConfig.SetDefault("searches", []string{})
}

func SaveRecentDownloads() {
	recentsConfig.Set("downloads", recentDownloads)
	recentsConfig.WriteConfigAs(recentsConfigPath)
}

func AddRecentDownload(info *scraper.DramaInfo) {
	recentDownloads[strings.ToLower(info.Name)] = []string{
		time.Now().Format(timeFormat),
		info.Name,
		info.SubURL,
	}
	SaveRecentDownloads()
}

// func AddRecentSearch(s string) {
// 	recentSearches = append(recentSearches, s)
// 	SaveRecents()
// }

type RecentEntry struct {
	key  string
	time int
}
type REList []RecentEntry

func (e REList) Len() int {
	return len(e)
}

func (e REList) Less(i, j int) bool {
	return e[i].time < e[j].time
}

func (e REList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func GetRecentDownloads() []scraper.DramaInfo {
	var dramas []scraper.DramaInfo
	var sorted []RecentEntry
	for key, props := range recentDownloads {
		res, err := time.Parse(timeFormat, props[0])
		if err != nil {
			panic(err)
		}
		sorted = append(sorted, RecentEntry{
			key:  key,
			time: int(time.Since(res).Seconds()),
		})
	}

	sort.Sort(REList(sorted))
	for _, ent := range sorted {
		key := ent.key
		props := recentDownloads[key]
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

// func GetRecentSearches() []string {
// 	return recentSearches
// }
