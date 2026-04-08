/*
Copyright © 2026 Mohamed Karim Benkirane <benkiranemedkarim@gmail.com>
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	providers   []string
	regionsFlag []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudping-go",
	Short: "A multi-cloud network latency benchmarking CLI",
	Long: `cloudping-go is a high-performance command-line tool designed to measure 
the true network latency between your machine and physical cloud data centers 
worldwide. It currently supports native endpoints for AWS, Google Cloud, and Azure.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&providers, "providers", "p", nil, "Filter by cloud provider(s) (comma-separated, e.g., aws,gcp)")
	rootCmd.PersistentFlags().StringSliceVarP(&regionsFlag, "regions", "r", nil, "Filter by region name(s) (comma-separated, e.g., us-east-1,europe-west9)")
}
