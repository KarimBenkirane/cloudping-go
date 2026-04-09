/*
Copyright © 2026 Mohamed Karim Benkirane <benkiranemedkarim@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/KarimBenkirane/cloudping-go/internal/pinger"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var pingCount int

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Measure HTTP-based Round-Trip Time (RTT) to cloud regions",
	Long: `Execute parallel latency tests against actual cloud data centers. 

This command dispatches highly-optimized, cache-busting HTTP requests directly 
to region-locked endpoints. By bypassing global edge networks, it reveals the 
true physical Round-Trip Time (RTT) to the requested servers. Results are 
automatically sorted from fastest to slowest.

Examples:
  cloudping-go ping
  cloudping-go ping --providers aws,gcp
  cloudping-go ping --regions us-east-1,europe-west9 -n 5`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if pingCount == 0 {
			return fmt.Errorf("Count cannot be 0!")
		}
		return nil
	},
	Run: runPing,
}

var sema = make(chan struct{}, 5) // semaphore for 5 workers
var wg sync.WaitGroup

func runPing(cmd *cobra.Command, args []string) {
	regions, err := pinger.FilterRegions(providers, regionsFlag)
	if err != nil {
		log.Fatal(err)
	}

	store := &pinger.ResultStore{}

	bar := progressbar.Default(int64(len(regions)), "Pinging servers...")

	for _, region := range regions {
		wg.Go(
			func() {
				sema <- struct{}{}
				result := pinger.PingRegion(region, pingCount)
				<-sema
				store.Add(result)
				bar.Add(1)
			})
	}
	wg.Wait()
	results := store.All()
	sort.Slice(results, func(i, j int) bool {
		if results[i].Status == "success" && results[j].Status != "success" {
			return true
		}
		if results[i].Status != "success" && results[j].Status == "success" {
			return false
		}
		return results[i].Latency < results[j].Latency
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
	rootCmd.SilenceUsage = true
	pingCmd.Flags().IntVarP(&pingCount, "count", "n", 3, "Number of times to ping each region (results are averaged)")
}
