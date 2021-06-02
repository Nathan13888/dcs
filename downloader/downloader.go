package downloader

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"dcs/prompt"

	"github.com/Nathan13888/m3u8/dl"
	"github.com/cavaliercoder/grab"
	"github.com/cheggaaa/pb/v3"
)

// Get - Download something
func Get(info DownloadInfo, prop DownloadProperties) error {
	var writer io.Writer = os.Stdout
	if info.Logger != nil {
		writer = info.Logger
	}

	overwrite := prop.Overwrite
	interactive := prop.Interactive
	ignorem3u8 := prop.IgnoreM3U8

	start := time.Now()
	var err error

	pathInfo := GetPath(info)
	dir := pathInfo.Dir
	path := pathInfo.Path
	partPath := path + ".part"

	// Create paths and directories
	LogInfo(writer, "Creating path '%s'\n\n", dir)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	// Check if something is already downloaded
	if Lookup(path) {
		LogInfo(writer, "The desired file '%s' has been already downloaded...\n", path)
		if overwrite {
			LogInfo(writer, "\nOverwriting '%s' because of OVERWRITE flag...\n", path)
		} else if interactive {
			if prompt.Confirm("Would you like to overwrite file?") {
				LogInfo(writer, "\nRemoving '%s'\n\n", path)
				err = os.Remove(path)
				if err != nil {
					return err
				}
			} else {
				LogInfo(writer, "\nSKIPPING download for '%s'...\n\n", path)
				return nil
				// return fmt.Errorf("user chose to not overwrite existing file")
			}
		} else {
			return fmt.Errorf("the download location '%s' already contains the episode", path)
		}
	}
	// TODO: more advanced lookup and incorporate checksums

	setupTime := time.Since(start)
	downloadStart := time.Now()

	LogInfo(writer, "Downloading '%s' EPISODE %v (%v)'\n\n", info.Name, info.Num, info.Link)
	if strings.HasSuffix(info.Link, ".mp4") { // is MP4
		err = DownloadMP4(writer, info, path, partPath)
	} else if strings.HasSuffix(info.Link, ".m3u8") { // is M3U8
		if !ignorem3u8 || (interactive && prompt.Confirm("Would you like to download M3U8 still?")) { // m3u8 is not ignored OR interactive prompt is confirmed
			err = DownloadM3U8(writer, info.Link, path)
		} else { //m3u8 is ignored
			return errors.New("m3u8 is ignored")
		}
	} else { // unknown
		return errors.New("unsupported file ending from '" + info.Link + "'")
	}
	if err != nil {
		return err
	}

	downloadTime := time.Since(downloadStart)

	LogInfo(writer, "\nFinished downloading '%s' EPISODE %v\n\n", info.Name, info.Num)

	// TODO: Scrap Time
	LogInfo(writer, "* Setup Time      >> %v\n", setupTime)
	LogInfo(writer, "* Download Time   >> %v\n", downloadTime)

	return nil
}

func DownloadMP4(writer io.Writer, info DownloadInfo, path string, partPath string) error {
	var err error

	// Start downloading
	client := grab.NewClient()
	req, _ := grab.NewRequest(partPath, info.Link)
	res := client.Do(req)
	LogInfo(writer, "Response: %v\n\n", res.HTTPResponse.Status)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	bar := pb.Full.Start64(res.Size)
	bar.SetRefreshRate(500 * time.Millisecond)
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
	bar.SetWriter(writer)
	if err = bar.Err(); err != nil {
		return err
	}
	defer bar.Finish()

Loop:
	for {
		select {
		case <-t.C:
			bar.SetCurrent(res.BytesComplete())
			if info.ProgressUpdater != nil {
				info.ProgressUpdater(float64(res.BytesComplete()) / float64(res.Size))
			}

		case <-res.Done:
			if info.ProgressUpdater != nil {
				info.ProgressUpdater(1.00) // 100% done
			}
			break Loop
		}
	}

	if err = res.Err(); err != nil {
		return err
	}

	LogInfo(writer, "\nDownload completed. Renaming file to final destination.\n\n")
	err = os.Rename(partPath, path)
	if err != nil {
		return err
	}

	return nil
}

func DownloadM3U8(writer io.Writer, url string, p string) error {
	var err error

	streams := 4 // number of concurrent downloaders
	tmpPath := p + "_m3u8files"
	mergedFile := path.Join(tmpPath, "main.ts")

	downloader, err := dl.NewTask(tmpPath, url)
	if err != nil {
		return err
	}
	// TODO: redirect output
	if err := downloader.Start(streams); err != nil {
		return err
	}

	LogInfo(writer, "\nFinished merging M3U8 files. Renaming file to final destination.\n\n")
	err = os.Rename(mergedFile, p)
	if err != nil {
		return err
	}

	LogInfo(writer, "Cleaning up... Removing '%s'\n\n", tmpPath)
	err = os.RemoveAll(tmpPath)
	if err != nil {
		return err
	}

	return err
}
