package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse content on remote.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Browsing remote...")
		u := GetRemoteURL("content/") + "/"
		open(u)
	},
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func init() {
	serviceCmd.AddCommand(browseCmd)
}
