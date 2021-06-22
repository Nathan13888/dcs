package downloader

import (
	"dcs/scraper"
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type EpisodeInfo struct {
	Num      float64 `json:"num"`
	Date     string  `json:"downloadTime,omitempty"`
	Hash     string  `json:"hash"`
	Attempts int     `json:"attempts"`
}

type CollectionInfo struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	FirstDT     string        `json:"firstDownloadTime"`
	LastDT      string        `json:"lastDownloadTime"`
	LastUpdated string        `json:"LastUpdated"`
	Episodes    []EpisodeInfo `json:"episodes"`
}

func WriteCollectionInfo(col string, info CollectionInfo) error {
	file, err := os.OpenFile(GetInfoPath(col), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	encoder.Encode(info)
	return nil
}

func ReadCollectionInfo(col string) (CollectionInfo, error) {
	var ret CollectionInfo
	file, err := os.OpenFile(GetInfoPath(col), os.O_RDONLY, 0644)
	if err != nil {
		return ret, err
	}
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&ret)
	if err != nil {
		fmt.Println(err)
	}
	return ret, nil
}

func UpdateDownload(col string, num float64, epName string) {

}

func GenerateInfo(col scraper.DramaInfo) {

}

func CalculateHash(ep scraper.EpisodeInfo) {
	// path := GetPath()
}

func GetInfoPath(col string) string {
	return path.Join(scraper.EscapeName(col), "info.txt")
}
