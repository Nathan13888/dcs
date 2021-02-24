package scraper

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"dcs/config"

	"github.com/cavaliercoder/grab"
	"github.com/cheggaaa/pb/v3"
)

// DownloadInfo - Information you need to Download
type DownloadInfo struct {
	Link string
	Name string
	Num  int
}

// Download - Download something
func Download(info DownloadInfo) error {
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

	// TODO: look up if target file exists and show prompt; accept flags

	// Start downloading
	client := grab.NewClient()
	req, _ := grab.NewRequest(partPath, info.Link)
	fmt.Printf("Downloading '%v'\n\n", req.URL())
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

	fmt.Printf("\nDownload completed. Renaming file to final name.\n\n")
	err = os.Rename(path+".part", path)
	if err != nil {
		return err
	}

	downloadTime := time.Since(res.Start)
	// TODO: Scrap Time
	fmt.Printf("* Setup Time      >> %v\n", setupTime)
	fmt.Printf("* Download Time   >> %v\n", downloadTime)

	return nil
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
