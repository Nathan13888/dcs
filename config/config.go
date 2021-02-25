package config

import (
	"path"

	"github.com/spf13/viper"
)

func SetDefaults(home string) {
	viper.SetDefault("DownloadPath", path.Join(home, "Downloads", "DCS"))
	viper.SetDefault("DaemonHost", "localhost")
	viper.SetDefault("DaemonPort", 6969)
}

func DownloadPath() string {
	return viper.GetString("DownloadPath")
}

func DaemonURL() (string, int) {
	return viper.GetString("DaemonHost"), viper.GetInt("DaemonPort")
}
