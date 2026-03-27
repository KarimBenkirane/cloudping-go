/*
Copyright © 2026 Mohamed Karim Benkirane <benkiranemedkarim@gmail.com>
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	providers []string
	codes     []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloudping-go",
	Short: "A multi-cloud latency testing tool",
	Long: `cloudping-go is a CLI utility written in Go 
that allows users to benchmark their network connection against 
the world's leading cloud infrastructures.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&providers, "providers", "p", nil, "Providers (eg. aws, gcp, azure)")
	rootCmd.PersistentFlags().StringSliceVarP(&codes, "codes", "c", nil, "Codes (eg. us-east-1)")
}
