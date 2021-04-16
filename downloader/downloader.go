package downloader

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"dcs/config"

	"github.com/Nathan13888/m3u8/dl"
	"github.com/cavaliercoder/grab"
	"github.com/cheggaaa/pb/v3"
)

// DownloadInfo - Information you need to Download
type DownloadInfo struct {
	Link string
	Name string
	Num  int
}

// Get - Download something
func Get(info DownloadInfo) error {
	start := time.Now()
	var err error

	pathInfo := GetPath(info)
	dir := pathInfo.Dir
	path := pathInfo.Path
	partPath := path + ".part"

	// Create paths and directories
	fmt.Printf("Creating path '%s'\n\n", dir)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// Check if something is already downloaded
	if Lookup(path) {
		fmt.Printf("The desired file '%s' has been already downloaded...\n", path)
		return errors.New("the download location '" + path + "' already contains the episode")
	}
	// TODO: more advanced lookup and incorporate checksums

	setupTime := time.Since(start)
	downloadStart := time.Now()

	// TODO: look up if target file exists and show prompt; accept flags

	fmt.Printf("Downloading '%v'\n\n", info.Link)
	if strings.HasSuffix(info.Link, ".mp4") {
		err = DownloadMP4(info, path, partPath)
	} else if strings.HasSuffix(info.Link, ".m3u8") {
		err = DownloadM3U8(info.Link, path)
	} else {
		return errors.New("unsupported file ending from '" + info.Link + "'")
	}
	if err != nil {
		return err
	}

	downloadTime := time.Since(downloadStart)
	// TODO: Scrap Time
	fmt.Printf("* Setup Time      >> %v\n", setupTime)
	fmt.Printf("* Download Time   >> %v\n", downloadTime)

	return nil
}

func DownloadMP4(info DownloadInfo, path string, partPath string) error {
	var err error
	// Start downloading
	client := grab.NewClient()
	req, _ := grab.NewRequest(partPath, info.Link)
	res := client.Do(req)
	fmt.Printf("Response: %v\n\n", res.HTTPResponse.Status)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	bar := pb.Full.Start64(res.Size)
	bar.SetRefreshRate(500 * time.Millisecond)
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
	if err = bar.Err(); err != nil {
		return err
	}
	defer bar.Finish()

Loop:
	for {
		select {
		case <-t.C:
			bar.SetCurrent(res.BytesComplete())

		case <-res.Done:
			break Loop
		}
	}

	if err = res.Err(); err != nil {
		return err
	}

	fmt.Printf("\nDownload completed. Renaming file to final destination.\n\n")
	err = os.Rename(partPath, path)
	if err != nil {
		return err
	}

	return nil
}

func DownloadM3U8(url string, p string) error {
	var err error

	streams := 4 // number of concurrent downloaders
	tmpPath := p + "_m3u8files"
	mergedFile := path.Join(tmpPath, "merged.ts")

	downloader, err := dl.NewTask(tmpPath, url)
	if err != nil {
		return err
	}
	if err := downloader.Start(streams); err != nil {
		return err
	}

	fmt.Printf("\nFinished merging M3U8 files. Renaming file to final destination.\n\n")
	err = os.Rename(mergedFile, p)
	if err != nil {
		return err
	}

	fmt.Printf("Cleaning up... Removing '%s'\n\n", tmpPath)
	err = os.RemoveAll(tmpPath)
	if err != nil {
		return err
	}

	return err
}

// Lookup - Check if downloaded content exists
func Lookup(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// PathInfo - Path information about downloaded content
type PathInfo struct {
	Folder  string
	Episode string
	Dir     string
	Path    string
}

// GetPath - Get PathInfo based on DownloadInfo
func GetPath(info DownloadInfo) PathInfo {
	folder := info.Name
	episode := fmt.Sprintf("ep%d.mp4", info.Num)
	dir := path.Join(config.DownloadPath(), folder)
	path := path.Join(dir, episode)
	return PathInfo{
		Folder:  folder,
		Episode: episode,
		Dir:     dir,
		Path:    path,
	}
}
