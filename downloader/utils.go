package downloader

import (
	"dcs/config"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func LogInfo(writer io.Writer, format string, a ...interface{}) {
	fmt.Fprintf(writer, format, a...)
}

func GetEpisodeNames(collection string) (int, []string, error) {
	cnt := 0
	var validEpisodes []string

	path := path.Join(config.DownloadPath(), collection)
	_, err := os.Stat(path)
	if err != nil {
		return -1, validEpisodes, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return -1, validEpisodes, err
	}

	mediaExtensions := []string{".mp4"}
	for _, file := range files {
		for _, ext := range mediaExtensions {
			if strings.EqualFold(ext, filepath.Ext(file.Name())) {
				// fmt.Printf("FOUND: %s\n", file.Name())
				validEpisodes = append(validEpisodes, file.Name())
				cnt++
				break
			}
		}
	}
	return cnt, validEpisodes, nil
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

func DownloadedDramas() (int, error) {
	cnt := 0
	path := config.DownloadPath()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return -1, err
	}
	for _, info := range files {
		if info.IsDir() {
			cnt++
		}
	}
	return cnt, nil
}

func DownloadedEpisodes() (int, error) {
	cnt := 0
	path := config.DownloadPath()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return -1, err
	}
	for _, info := range files {
		if info.IsDir() {
			c, _, err := GetEpisodeNames(info.Name())
			if err != nil {
				break // if statement at the end will process return
			}
			cnt += c
		}
	}
	if err != nil {
		return -1, err
	}
	return cnt, nil
}

func Size(rpath string) (int64, error) {
	info, err := os.Stat(path.Join(config.DownloadPath(), rpath))
	if err != nil {
		return -1, err
	}
	if info.IsDir() {
		return DirSize(rpath)
	}
	return info.Size(), nil
}

func DirSize(rpath string) (int64, error) {
	var size int64
	path := path.Join(config.DownloadPath(), rpath)
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	if err != nil {
		return -1, err
	}
	return size, nil
}
