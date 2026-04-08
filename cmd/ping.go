/*
Copyright © 2026 Mohamed Karim Benkirane <benkiranemedkarim@gmail.com>
*/
package cmd

import (
	"log"
	"sort"
	"sync"

	"github.com/KarimBenkirane/cloudping-go/internal/pinger"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var pingCount int64

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Measure HTTP latency to cloud regions",
	Long: `Execute latency tests against global cloud endpoints (AWS, GCP, Azure). 
This command measures the round-robin time (RTT) to the nearest entry points 
of various cloud providers, helping you identify the fastest region from 
your current network location.

Example:
  cloudping-go ping --providers aws --regions us-east-1,us-west-2`,
	Run: runPing,
}

var sema = make(chan struct{}, 5) // semaphore for 5 workers
var wg sync.WaitGroup

func runPing(cmd *cobra.Command, args []string) {
	regions, err := pinger.FilterRegions(providers, codes)
	if err != nil {
		log.Fatal(err)
	}

	results := make([]pinger.Result, 0)
	bar := progressbar.Default(int64(len(regions)), "Pinging servers...")

	for _, region := range regions {
		wg.Go(
			func() {
				sema <- struct{}{}
				result := pinger.PingRegion(region, pingCount)
				<-sema
				results = append(results, result)
				bar.Add(1)
			})
	}
	wg.Wait()
	sort.Slice(results, func(i, j int) bool {
		return results[i].Latency <= results[j].Latency
	})
	printTable(results)

}

func printTable(results []pinger.Result) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Provider", "Name", "Code", "Latency (ms)", "Status")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, result := range results {
		var latencyColor string
		if result.Latency > 0 && result.Latency < 80 {
			latencyColor = color.GreenString("%d", result.Latency) // Green
		} else if result.Latency >= 80 && result.Latency <= 200 {
			latencyColor = color.YellowString("%d", result.Latency) // Yellow
		} else {
			latencyColor = color.RedString("%d", result.Latency) // Red
		}
		tbl.AddRow(result.Region.Provider, result.Region.Name, result.Region.Code, latencyColor, result.Status)
	}

	tbl.Print()
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().Int64VarP(&pingCount, "count", "t", 3, "Define the amount of times to ping the server (the result will be the average of those pings)")
}
