package cmd

import (
	"dcs/server"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List of jobs",
	Long:  `List download jobs on the remote`,
	Run: func(cmd *cobra.Command, args []string) {
		u := GetRemoteURL("api/jobs")

		res, err := http.Get(u)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		var jobs server.JobsResponse
		decoder := json.NewDecoder(res.Body)
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&jobs)
		if err != nil {
			panic(err)
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Collection", "Episode", "Status", "Size"})
		var sum int64
		for i, job := range jobs.Jobs {
			size := jobs.Sizes[i]
			sum += size
			row := []string{
				job.ID,
				job.Req.DInfo.Name,
				fmt.Sprintf("%v", job.Req.DInfo.Num),
				string(job.Status),
				fmt.Sprintf("%.2f GB", float64(size)/math.Pow(1024, 3)),
			}
			table.Append(row)
		}
		table.SetFooter([]string{"", "", "", "Total Size", fmt.Sprintf("%.1f GB", float64(sum)/math.Pow(1024, 3))})

		table.SetBorder(false)
		table.SetAutoWrapText(true)
		table.SetAutoMergeCellsByColumnIndex([]int{1}) // merge dramas
		table.Render()
	},
}

func init() {
	serviceCmd.AddCommand(listCmd)
}
