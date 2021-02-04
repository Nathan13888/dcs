package scraper

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
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
	fmt.Printf("Creating path '%s.part'\n\n", path)
	err = os.MkdirAll(dir+".part", 0755)
	if err != nil {
		return err
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Start downloading
	client := grab.NewClient()
	req, _ := grab.NewRequest(path+".part", info.Link)
	fmt.Printf("Downloading '%v'\n\n", req.URL())
	res := client.Do(req)
	fmt.Printf("  %v\n", res.HTTPResponse.Status)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %s / %s (%.2f%%)\n",
				convertSize(res.BytesComplete()),
				convertSize(res.Size),
				100*res.Progress())

		case <-res.Done:
			break Loop
		}
	}

	// check for errors
	if err := res.Err(); err != nil {
		return err
	}

	fmt.Printf("Download completed. Renaming file to final name.\n\n")
	err = os.Rename(path+".part", path)
	if err != nil {
		return err
	}

	return nil
}

func convertSize(bytes int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	pow := int64(math.Log(float64(bytes)) / math.Log(1024))
	return fmt.Sprintf("%.2f %s",
		float64(bytes)/math.Pow(1024, float64(pow)),
		units[pow])
}

func Lookup() bool {
	return false
}
