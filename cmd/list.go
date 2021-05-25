package cmd

import (
	"dcs/server"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List of jobs",
	Long:  `List download jobs on the remote`,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := Request("GET", "api/jobs")
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

		if len(jobs.Jobs) == 0 { // no jobs found
			fmt.Println("No jobs found...")
			return
		}
		if len(jobs.Jobs) != len(jobs.Sizes) {
			panic(fmt.Errorf("invalid response received"))
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Collection", "Episode", "Status", "Progress", "Scheduled", "Download Time", "Date", "Size"})
		var sum int64
		var totalProgress float64
		var totalDur time.Duration
		for i, job := range jobs.Jobs {
			size := jobs.Sizes[i]
			sum += size
			totalProgress += job.Progress.Completion

			dt := "unknown"
			if !job.Progress.StartTime.IsZero() {
				var t time.Time
				if job.Progress.EndTime.IsZero() {
					t = time.Now()
				} else {
					t = job.Progress.EndTime
				}
				dur := t.Sub(job.Progress.StartTime)
				totalDur += dur
				dt = dur.Round(time.Second).String()
			}

			row := []string{
				job.ID,
				job.Req.DInfo.Name,
				fmt.Sprintf("%v", job.Req.DInfo.Num),
				string(job.Progress.Status),
				fmt.Sprintf("%.2f %%", job.Progress.Completion),
				job.Schedule.Format(time.RFC822),
				dt,
				job.Date.Format(time.RFC822),
				fmt.Sprintf("%.2f GB", float64(size)/math.Pow(1024, 3)),
			}
			table.Append(row)
		}
		table.SetFooter([]string{"", "", "",
			"Total Progress", fmt.Sprintf("%.1f %%",
				totalProgress/float64(len(jobs.Jobs))),
			"Total DT", totalDur.Round(time.Millisecond).String(),
			"Total Size", fmt.Sprintf("%.1f GB",
				float64(sum)/math.Pow(1024, 3)),
		})

		table.SetBorder(false)
		table.SetAutoWrapText(true)
		table.SetAutoMergeCellsByColumnIndex([]int{1}) // merge dramas
		table.Render()
	},
}

func init() {
	serviceCmd.AddCommand(listCmd)
}
