package scraper

import (
	"fmt"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/cheggaaa/pb/v3"
	"github.com/mitchellh/go-homedir"
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

	// Paths
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	folder := info.Name
	episode := fmt.Sprintf("ep%d.mp4", info.Num)
	// TODO: config download location
	dir := fmt.Sprintf("%s/Downloads/DCS/%s", home, folder)
	path := dir + "/" + episode

	// Create paths and directories
	fmt.Printf("Creating path '%s'\n\n", dir)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	setupTime := time.Since(start)

	// TODO: look up if target file exists and show prompt; accept flags

	// Start downloading
	client := grab.NewClient()
	req, _ := grab.NewRequest(path+".part", info.Link)
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

func Lookup() bool {
	return false
}
