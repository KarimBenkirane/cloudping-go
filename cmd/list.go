/*
Copyright © 2026 Mohamed Karim Benkirane <benkiranemedkarim@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/KarimBenkirane/cloudping-go/internal/pinger"

	"github.com/spf13/cobra"

	"github.com/fatih/color"

	"github.com/rodaine/table"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Display all supported cloud providers and regions",
	Long: `Print a comprehensive, formatted table of all active AWS, GCP, and Azure 
regions available for latency testing. 

Use this command to quickly find the exact Region strings needed for the ping 
command's '--regions' filter, or to verify the endpoint URLs being tested.

Examples:
  cloudping-go list --providers azure
  cloudping-go list --regions us-east-1,eu-west-1`,
	Run: runList,
}

func runList(cmd *cobra.Command, args []string) {
	regions, err := pinger.FilterRegions(providers, regionsFlag)
	if err != nil {
		log.Fatal(err)
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Provider", "Name", "Code", "Url")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, region := range regions {
		tbl.AddRow(region.Provider, region.Name, region.Code, region.Url)
	}

	tbl.Print()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
