/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"dcs/server"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		job, err := cmd.Flags().GetString("job")
		if err != nil {
			panic(err)
		}
		var rp string = "api/log"
		if len(job) > 0 { // display logs about job (if it exists)
			rp += "/job/" + job
		} else {
			rp += "/server"
		}
		// display logs of current remote session
		res, err := Request("GET", rp)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		fmt.Println("Status response:", res.StatusCode)

		logObj := server.LogLookupResponse{}
		decoder := json.NewDecoder(res.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&logObj)
		if err != nil {
			panic(err)
		}
		displayLog(logObj.Log)
	},
}

// TODO: log streaming
// TODO: better log navigation
func displayLog(log []string) {
	fmt.Println("Displaying logs...")

	cw := zerolog.ConsoleWriter{Out: os.Stdout}
	for _, line := range log {
		cw.Write([]byte(line))
	}
}

func init() {
	serviceCmd.AddCommand(logsCmd)

	logsCmd.Flags().StringP("job", "j", "", "Display logs of a job instead")
}
