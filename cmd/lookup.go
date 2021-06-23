package cmd

import (
	"dcs/config"
	"dcs/prompt"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// lookupCmd represents the lookup command
var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Check downloaded content",
	Long: `Lookup information about downloaded content.

	USAGE: lookup`,
	Run: func(cmd *cobra.Command, args []string) {
		// All the folders of the dramas
		fmt.Printf("Listing all the folders\n\n")
		var files []os.FileInfo
		files = getFiles(config.DownloadPath())
		// ! Note, if there are no folders in the directory, then sucks for you :P
		var res os.FileInfo
		res, _ = prompt.DirSelect("Select a folder", files, true)
		fmt.Printf("Selected '%s'\n", res.Name())

		path := path.Join(config.DownloadPath(), res.Name())
		fmt.Println(path)
		files = getFiles(path)
		res, _ = prompt.DirSelect("Select an episode", files, false)
		fmt.Printf("Selected '%s'\n", res.Name())
	},
}

func getFiles(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	return files
}

func init() {
	rootCmd.AddCommand(lookupCmd)
}
