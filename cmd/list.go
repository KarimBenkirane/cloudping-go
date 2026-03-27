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
	Short: "List all supported cloud regions",
	Long: `Display a comprehensive list of all AWS, GCP, and Azure regions 
available for testing. You can filter the list by provider or region code 
to verify which endpoints are currently configured in the tool.

Example:
  cloudping-go list --providers gcp`,
	Run: runList,
}

func runList(cmd *cobra.Command, args []string) {
	regions, err := pinger.FilterRegions(providers, codes)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
