package downloader

import (
	"dcs/config"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// TODO: improve results and make caller handle results
func DisplayEpisodes(name string) int {
	cnt := 0
	path := path.Join(config.DownloadPath(), name)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return cnt
	} else if !os.IsExist(err) && err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(path)
	// might not be cool to crash because someone gave this shit to read
	if err != nil {
		panic(err)
	}

	mediaExtensions := []string{".mp4"}
	for _, file := range files {
		for _, ext := range mediaExtensions {
			if strings.EqualFold(ext, filepath.Ext(file.Name())) {
				fmt.Printf("FOUND: %s\n", file.Name())
				cnt++
				break
			}
		}
	}
	return cnt
}

// Lookup - Check if downloaded content exists
func Lookup(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetPath - Get PathInfo based on DownloadInfo
func GetPath(info DownloadInfo) PathInfo {
	folder := info.Name
	episode := fmt.Sprintf("ep%v.mp4", info.Num)
	dir := path.Join(config.DownloadPath(), folder)
	path := path.Join(dir, episode)
	return PathInfo{
		Folder:  folder,
		Episode: episode,
		Dir:     dir,
		Path:    path,
	}
}

func Size(rpath string) int {
	return 0
}
