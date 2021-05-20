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

const timeFormat = "Jan-02-06_15:04:05"

var recentsConfigPath string

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
	ref  *[]string
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
	L := len(recentDownloads)
	dramas := make([]scraper.DramaInfo, L)
	sorted := make([]RecentEntry, L)
	i := 0
	for _, props := range recentDownloads {
		res, err := time.Parse(timeFormat, props[0])
		copy := props
		if err != nil {
			panic(err)
		}
		sorted[i] = RecentEntry{
			ref:  &copy,
			time: int(time.Since(res).Seconds()),
		}
		i++
	}

	sort.Sort(REList(sorted))
	j := 0
	for _, ent := range sorted {
		props := *ent.ref
		if !(len(props) >= 2) {
			panic(fmt.Errorf("invalid properties for recent downloads entry"))
		}
		name := props[1]
		subURL := props[2]

		dramas[j] = scraper.DramaInfo{
			FullURL: scraper.URL + subURL,
			SubURL:  subURL,
			Domain:  scraper.URL,
			Name:    name,
		}
		j++
	}

	return dramas
}

// func GetRecentSearches() []string {
// 	return recentSearches
// }
