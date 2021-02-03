package scraper

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/mitchellh/go-homedir"
)

type DownloadInfo struct {
	Link string
	Name string
	Num  int
}

// Download - Download something
func Download(info DownloadInfo) error {
	res, err := http.Get(info.Link)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	home, err := homedir.Dir()

	if err != nil {
		return err
	}

	folder := info.Name
	// TODO: add .part extension
	episode := fmt.Sprintf("ep%d.mp4", info.Num)

	// TODO: config download location
	dir := fmt.Sprintf("%s/Downloads/DCS/%s", home, folder)
	path := dir + "/" + episode

	fmt.Printf("Creating path '%s'\n\n", path)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}

	defer out.Close()

	fmt.Printf("Saving file to '%s'\n\n", path)
	_, err = io.Copy(out, res.Body)

	if err != nil {
		return err
	}

	// TODO: check downloaded file size
	// TODO: check downloaded file playable

	return nil
}

func Lookup() bool {
	return false
}
