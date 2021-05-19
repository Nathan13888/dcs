package config

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var IS_SERVER = false
var saved_home string

func GetHome() string {
	if len(saved_home) > 0 { // return saved home directory
		return saved_home
	}
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

func GetConfigHome() string {
	return path.Join(GetHome(), ".dcs")
}

func Configure() {
	// Search config in home directory with name ".dcs" (without extension).
	viper.SetConfigName("config.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(GetConfigHome())
	viper.AddConfigPath("/etc/dcs/")
	viper.AddConfigPath(".")
	SetDefaults()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		fmt.Println("CONFIG FILE NOT FOUND, CREATING DEFAULT CONFIG")
		os.MkdirAll(GetConfigHome(), 0755)
		viper.SafeWriteConfigAs(path.Join(GetConfigHome(), "config.json"))
	} else {
		panic(err)
	}
}

func SetDefaults() {
	viper.SetDefault("DownloadPath", path.Join(GetHome(), "Downloads", "DCS"))
	viper.SetDefault("DaemonHost", "localhost")
	viper.SetDefault("DaemonPort", 6969)
	viper.SetDefault("SQLiteFile", path.Join(GetConfigHome()))
}

func DownloadPath() string {
	return viper.GetString("DownloadPath")
}

func DaemonURL() (string, int) {
	return viper.GetString("DaemonHost"), viper.GetInt("DaemonPort")
}
