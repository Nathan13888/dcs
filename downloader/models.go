package downloader

// DownloadInfo - Information you need to Download
type DownloadInfo struct {
	Link string
	Name string
	Num  float64
}

type DownloadProperties struct {
	Overwrite   bool
	Interactive bool
	IgnoreM3U8  bool
	Remote      bool
}

// PathInfo - Path information about downloaded content
type PathInfo struct {
	Folder  string
	Episode string
	Dir     string
	Path    string
}
