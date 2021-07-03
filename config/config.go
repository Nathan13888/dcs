package config

import (
	"fmt"
	"os"
	"path"
	"strings"

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
	if IS_DEV() {
		viper.SetConfigName("config.dev.json")
	} else {
		viper.SetConfigName("config.json")
	}
	viper.SetConfigType("json")
	viper.AutomaticEnv()
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
	viper.SetDefault("DownloadMethod", DEFAULTMETHOD) // refer to DMethod
	viper.SetDefault("DaemonHost", "localhost")
	viper.SetDefault("DaemonPort", 6969)
	viper.SetDefault("EnableFileServer", false)
	viper.SetDefault("DownloadLimit", 2)
	viper.SetDefault("DSN", path.Join(GetConfigHome(), "dcs.db"))
	// viper.SetDefault("DSN",
	// 	"host=localhost user=dcs password=dcspassword dbname=dcs port=9920 sslmode=disable TimeZone=America/Toronto")
}

func DownloadPath() string {
	return viper.GetString("DownloadPath")
}

type DMethod string

const (
	AjaxMethod    DMethod = "ajax"
	LDMethod      DMethod = "ld" // lookup download
	ManualMethod  DMethod = "manual"
	DEFAULTMETHOD DMethod = ManualMethod
)

func DownloadMethod() DMethod {
	s := strings.ToLower(viper.GetString("DownloadMethod"))
	// methods := []DMethod{
	// 	AjaxMethod,
	// 	LDMethod,
	// 	ManualMethod,
	// }
	// found := false
	// for _, m := range methods {
	// 	if m == DMethod(strings.ToLower(s)) {
	// 		found = true
	// 		continue
	// 	}
	// }
	// if !found {
	// 	return DEFAULTMETHOD
	// }
	return DMethod(s)
}

func DaemonURL() (string, int) {
	return viper.GetString("DaemonHost"), viper.GetInt("DaemonPort")
}

func EnableFileServer() bool {
	return viper.GetBool("EnableFileServer")
}

func DownloadLimit() int {
	return viper.GetInt("DownloadLimit")
}

func DSN() string {
	return viper.GetString("DSN")
}
