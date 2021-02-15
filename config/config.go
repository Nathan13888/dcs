package config

import "github.com/spf13/viper"

func DownloadPath() string {
	return viper.GetString("DownloadPath")
}
