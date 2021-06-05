package downloader

import "io"

// DownloadInfo - Information you need to Download
type DownloadInfo struct {
	Link            string        `json:"link"`
	Name            string        `json:"name"`
	Num             float64       `json:"num"`
	Logger          io.Writer     `gorm:"-" json:"-"`
	ProgressUpdater func(float64) `gorm:"-" json:"-"`
}

type DownloadProperties struct {
	Overwrite   bool `json:"overwrite"`
	Interactive bool `json:"interactive"`
	IgnoreM3U8  bool `json:"ignoreM3U8"`
	Remote      bool `json:"remote"`
}

// PathInfo - Path information about downloaded content
type PathInfo struct {
	Folder  string
	Episode string
	Dir     string
	Path    string
}
