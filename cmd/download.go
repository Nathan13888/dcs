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
	"os"
	"strings"

	"github.com/manifoldco/promptui"
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
		noRecent, err := cmd.Flags().GetBool("no-recent")
		if err != nil {
			panic(err)
		}
		bulkMode, err := cmd.Flags().GetBool("bulk")
		if err != nil {
			panic(err)
		}
		tryout, err := cmd.Flags().GetBool("tryout")
		if err != nil {
			panic(err)
		}
		manual, err := cmd.Flags().GetBool("manual")
		if err != nil {
			panic(err)
		}

		prop := downloader.DownloadProperties{
			Overwrite:   overwrite,
			Interactive: !interactive,
			IgnoreM3U8:  !ignorem3u8,
			Remote:      remote,
			ManualMode:  manual,
		}

		if len(args) == 1 && scraper.IsLink(args[0]) {
			download(args[0], prop)
		} else if len(args) >= 2 { // for first search and episode range
			link := scraper.FirstSearch(scraper.JoinArgs(args[:len(args)-1]))
			er := scraper.GetRange(args[len(args)-1]) // select last arg as range
			downloadRange(scraper.GetEpisodesByLink(link), er, prop)
		} else { // interactive prompted download(s)
			type Target struct {
				eps []scraper.EpisodeInfo
				er  []float64
			}
			var targets []Target
			for {
				var drama scraper.DramaInfo
				if !noRecent && !tryout {
					drama = *searchRecent(remote)
					// } else if enterLink {
					// drama=scraper.
				} else {
					drama = *searchDrama()
				}
				updateRecent(&drama, remote)

				cnt, epInfo, csize := lookupDownloadedEpisodes(drama, remote)
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

				var episodeRange []float64
				if tryout {
					episodeRange = []float64{1.0}
				} else {
					res, err := prompt.String("Episode Range")
					if err == promptui.ErrInterrupt {
						os.Exit(0)
					} else if err != nil {
						panic(err)
					}
					episodeRange = scraper.GetRange(strings.TrimSpace(res))
				}
				targets = append(targets, Target{
					eps: episodes,
					er:  episodeRange,
				})
				if !(bulkMode || tryout) { // do not prompt again
					break
				}
				if prompt.Confirm("Would you like to start downloading?") {
					break
				}
			}
			for _, tar := range targets {
				downloadRange(tar.eps, tar.er, prop)
			}
		}
	},
}

func downloadRange(episodes []scraper.EpisodeInfo, erange []float64, prop downloader.DownloadProperties) {
	fmt.Printf("Attemping to download these episodes: %v\n\n", erange)

	for i := len(episodes) - 1; i >= 0; i-- {
		e := episodes[i]
		url := e.Link
		for x := 0; x < len(erange); x++ { // erange should be sorted --> better efficiency
			if e.Number == erange[x] && len(url) > 0 {
				erange[x] = math.MaxFloat64
				download(scraper.URL+url, prop)
				break
			}
		}
	}
	for _, er := range erange {
		if er != math.MaxFloat64 {
			fmt.Printf("Episode %v was not available.\n", er)
		}
	}
}

func lookupDownloadedEpisodes(drama scraper.DramaInfo, remote bool) (int, []string, int64) {
	var cnt int
	var epInfo []string
	var csize int64
	var err error

	if remote {
		var obj server.CollectionLookupResponse
		res, err := Request("GET", "api/lookup/collection/"+drama.Name)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()

		code := res.StatusCode
		if code != http.StatusOK {
			return 0, []string{}, 0
		}

		decoder := json.NewDecoder(res.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&obj)
		if err != nil {
			fmt.Println(err)
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

	return cnt, epInfo, csize
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
	if err == promptui.ErrInterrupt {
		os.Exit(0)
	} else if err != nil {
		panic(err)
	}

	if *res == searchItem {
		fmt.Println("Searching for drama instead.")
		return searchDrama()
	}
	return res
}

func searchDrama() *scraper.DramaInfo {
	var drama *scraper.DramaInfo
	res, err := prompt.String("Search")
	if err == promptui.ErrInterrupt {
		os.Exit(0)
	} else if err != nil {
		panic(err)
	}
	queries := scraper.Search(res)
	if len(queries) == 0 {
		fmt.Printf("Found no results found with '%s'.\n", res)
		return searchDrama()
	} else {
		resInfo, err := prompt.Drama(queries)
		if err == promptui.ErrInterrupt {
			os.Exit(0)
		} else if err != nil {
			panic(err)
		}
		drama = resInfo
		//TODO: more rigorous checking
	}
	return drama
}

// TODO: update error handling and prompts
func download(episode string, prop downloader.DownloadProperties) {
	var dinfo downloader.DownloadInfo
	url := GetRemoteURL("api/download")

	fmt.Printf("Attemping to download from '%s'\n\n", episode)
	if prop.ManualMode {
		name, episodeNum, streaming := scraper.GetInfo(episode)
		fmt.Printf("\nFOUND STREAMING LINK: `%s`\n", streaming)
		manualLink, err := prompt.String(fmt.Sprintf("Enter link for %s #%v", name, episodeNum))
		if err == promptui.ErrAbort {
			os.Exit(0)
		}
		fmt.Printf("\nEntered MANUAL link: `%s`\n\n", manualLink)
		dinfo = downloader.DownloadInfo{
			Link: strings.Trim(manualLink, " \n"),
			Name: scraper.EscapeName(name),
			Num:  episodeNum,
		}
	} else {
		ajax := scraper.GetAjax(episode)
		if ajax.Found || (prop.Interactive && prompt.Confirm("Ajax not found. Would you like to proceed downloading?")) {
			fmt.Printf("Found AJAX endpoint '%s'\n\n", ajax.Ajax)
			link := scraper.ScrapeAjax(ajax)
			fmt.Printf("Found '%s'\n\n", link)
			// TODO: prompt confirm download
			dinfo = downloader.DownloadInfo{
				Link: link,
				Name: scraper.EscapeName(ajax.Name),
				Num:  ajax.Num,
			}
		} else {
			panic(fmt.Errorf("found bad AJAX: %v", ajax))
		}
	}
	if prop.Remote {
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
		fmt.Printf("Job Status: %s\n\n", job.Progress.Status)
		fmt.Printf("Sent job for %s EPISODE %v\n\n\n", job.Req.DInfo.Name, job.Req.DInfo.Num)
	} else {
		fmt.Println("Downloading...")
		err := downloader.Get(dinfo, prop)
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().BoolP("no-recent", "n", false, "Do not display recently downloaded dramas")
	downloadCmd.Flags().BoolP("overwrite", "o", false, "Overwrite if episode exists")
	downloadCmd.Flags().BoolP("no-interactive", "i", false, "Prompt to overwrite episode; important for automated download")
	downloadCmd.Flags().BoolP("dont-ignore-m3u8", "m", false, "Download M3U8 files")
	downloadCmd.Flags().BoolP("bulk", "b", false, "bulk mode")
	downloadCmd.Flags().BoolP("tryout", "t", false, "tryout a drama (just assume ep 1)")
	downloadCmd.Flags().BoolP("manual", "M", false, "manual mode (use own download link)")
}
