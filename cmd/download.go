package cmd

import (
	"bytes"
	"dcs/config"
	"dcs/downloader"
	"dcs/prompt"
	"dcs/scraper"
	"dcs/server"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download an episode or episodes of a drama",
	Long: `Download anything from DCS that you want.

	USAGE: download  -->  (for interactive prompt)
	USAGE: download <link to episode>
	USAGE: download <name of drama> <episode range>`,
	Aliases: []string{
		"down", "d",
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: sanitize arguments
		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			panic(err)
		}
		interactive, err := cmd.Flags().GetBool("no-interactive")
		if err != nil {
			panic(err)
		}
		ignorem3u8, err := cmd.Flags().GetBool("dont-ignore-m3u8")
		if err != nil {
			panic(err)
		}
		remote, err := cmd.Flags().GetBool("remote")
		if err != nil {
			panic(err)
		}
		if remote {
			if !scraper.Ping(config.DaemonURL()) {
				panic(fmt.Errorf("remote server NOT online"))
			}
		}

		prop := downloader.DownloadProperties{
			Overwrite:   overwrite,
			Interactive: !interactive,
			IgnoreM3U8:  !ignorem3u8,
			Remote:      remote,
		}

		if len(args) == 1 && scraper.IsLink(args[0]) {
			download(args[0], prop)
		} else {
			var link string
			var episodeRange []float64
			if len(args) == 0 {
				var res string
				var err error
				var drama scraper.DramaInfo
				showRecent, err := cmd.Flags().GetBool("no-recent")
				if err != nil {
					panic(err)
				}
				if !showRecent {
					drama = *searchRecent(remote)
				} else {
					drama = *searchDrama()
				}
				link = drama.FullURL

				updateRecent(&drama, remote)

				var cnt int
				var epInfo []string
				var csize int64

				if remote {
					var obj server.CollectionLookupResponse
					res, err := Request("GET", "api/lookup/collection/"+drama.Name)
					if err != nil {
						panic(err)
					}
					defer res.Body.Close()

					code := res.StatusCode
					if code != http.StatusOK {
						cnt = 0
						goto DisplayFoundEpisodes
					}

					decoder := json.NewDecoder(res.Body)
					decoder.DisallowUnknownFields()
					err = decoder.Decode(&obj)
					if err != nil {
						panic(err)
					}
					cnt = obj.NumOfEpisodes
					epInfo = obj.DownloadedEpisodes
					// err = obj.Error
					csize = obj.Size
				} else {
					cnt, epInfo, err = downloader.CollectionLookup(drama.Name)
					if err != nil {
						fmt.Println(err)
					}
					csize, err = downloader.DirSize(drama.Name)
					if err != nil {
						fmt.Println(err)
					}
				}
			DisplayFoundEpisodes:
				if cnt == 0 {
					fmt.Print("No episodes found.\n\n")
				} else {
					fmt.Printf("\nFound %d episodes:\n", cnt)
					for _, e := range epInfo {
						fmt.Printf("FOUND %s\n", e)
					}
					fmt.Printf("\nTotal size of collection: %.3f GB\n\n", float64(csize)/math.Pow(1024, 3))
				}

				episodes := scraper.GetEpisodes(drama)
				DisplayEpisodesInfo(episodes)

				res, err = prompt.String("Episode Range")
				if err != nil {
					panic(err)
				}
				episodeRange = scraper.GetRange(strings.TrimSpace(res))
			} else {
				link = scraper.FirstSearch(scraper.JoinArgs(args[:len(args)-1]))
				episodeRange = scraper.GetRange(args[len(args)-1])
			}
			if link != "" {
				// TODO: sanitize value of episode
				episodes := scraper.GetEpisodesByLink(link)
				fmt.Printf("Attemping to download these episodes: %v\n\n", episodeRange)
				for _, num := range episodeRange {
					fmt.Printf("Looking for episode %v...\n", num)
					var url string
					for _, e := range episodes {
						if e.Number == num {
							url = e.Link
						}
					}
					containsEpisode := false
					for _, ep := range episodes {
						if ep.Number == num {
							containsEpisode = true
							break
						}
					}
					if containsEpisode || url == "" {
						download(scraper.URL+url, prop)
					} else {
						fmt.Printf("Episode %v was not available", num)
					}
				}
			} else {
				fmt.Println("There has been a problem using your specified query")
				return
			}
		}
	},
}

func updateRecent(drama *scraper.DramaInfo, remote bool) {
	if remote {
		url := GetRemoteURL("api/recentdownload")

		json, err := json.Marshal(*drama)
		if err != nil {
			panic(err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		fmt.Printf("Received Status: %s\n\n", res.Status)
	} else {
		config.AddRecentDownload(drama)
	}
}

func searchRecent(remote bool) *scraper.DramaInfo {
	var recent []scraper.DramaInfo
	if remote {
		var obj []scraper.DramaInfo
		res, err := Request("GET", "api/recentdownloads")
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		decoder := json.NewDecoder(res.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&obj)
		if err != nil {
			panic(err)
		}

		recent = obj
	} else {
		recent = config.GetRecentDownloads()
	}
	if len(recent) == 0 {
		fmt.Println("No recent history. Searching instead.")
		return searchDrama()
	}
	searchItem := scraper.DramaInfo{
		Name:    "* SEARCH INSTEAD *",
		FullURL: "/link-to-no-where",
		SubURL:  "/link-to-no-where",
		Domain:  "notadomain.com",
	}

	res, err := prompt.Drama(append([]scraper.DramaInfo{searchItem}, recent...))
	if err != nil {
		panic(err)
	}

	// fmt.Println(res)
	// fmt.Println(searchItem)
	if *res == searchItem {
		fmt.Println("Searching for drama instead.")
		return searchDrama()
	}
	return res
}

func searchDrama() *scraper.DramaInfo {
	var drama *scraper.DramaInfo
	res, err := prompt.String("Search")
	if err != nil {
		panic(err)
	}
	queries := scraper.Search(res)
	if len(queries) == 0 {
		// TODO: don't PANIC
		panic("no queries were found")
	} else {
		resInfo, err := prompt.Drama(queries)
		if err != nil {
			panic(err)
		}
		drama = resInfo
		//TODO: more rigorous checking
	}
	return drama
}

func download(link string, prop downloader.DownloadProperties) {
	remote := prop.Remote

	url := GetRemoteURL("api/download")

	ajax := scraper.GetAjax(link)
	if ajax.Found || (prop.Interactive && prompt.Confirm("Ajax not found. Would you like to proceed downloading?")) {
		fmt.Printf("Attemping to download from '%s'\n\n", link)
		fmt.Printf("Found AJAX endpoint '%s'\n\n", ajax.Ajax)
		link := scraper.ScrapeAjax(ajax)
		fmt.Printf("Found '%s'\n\n", link)
		// TODO: prompt confirm download
		dinfo := downloader.DownloadInfo{
			Link: link,
			Name: ajax.Name,
			Num:  ajax.Num,
		}
		if remote {

			// TODO: change protocol

			jobinfo, err := json.Marshal(server.DownloadRequest{
				DInfo: dinfo,
				Props: prop,
			})
			if err != nil {
				panic(err)
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jobinfo))
			if err != nil {
				panic(err)
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			var job server.DownloadJob
			fmt.Printf("Received Status: %s\n", res.Status)
			decoder := json.NewDecoder(res.Body)
			decoder.DisallowUnknownFields()
			err = decoder.Decode(&job)
			if err != nil {
				panic(err)
			}

			fmt.Printf("Job ID:     %s\n", job.ID)
			fmt.Printf("Job Status: %s\n\n", job.Status)
			fmt.Printf("Job downloads %s EPISODE %v\n", job.Req.DInfo.Name, job.Req.DInfo.Num)
		} else {
			fmt.Println("Downloading...")
			err := downloader.Get(dinfo, prop)
			if err != nil {
				panic(err)
			}
		}
	} else {
		fmt.Print("FAILED to find episode...\n\n")
		fmt.Printf("Found AJAX: %v\n", ajax)
	}
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().BoolP("no-recent", "n", false, "Do not display recently downloaded dramas")
	downloadCmd.Flags().BoolP("overwrite", "o", false, "Overwrite if episode exists")
	downloadCmd.Flags().BoolP("no-interactive", "i", false, "Prompt to overwrite episode; important for automated download")
	downloadCmd.Flags().BoolP("dont-ignore-m3u8", "m", false, "Download M3U8 files")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
