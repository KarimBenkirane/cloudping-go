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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cloudping-go.yaml)")
	rootCmd.PersistentFlags().StringSliceVar(&providers, "providers", nil, "Providers (eg. aws, gcp, azure)")
	rootCmd.PersistentFlags().StringSliceVar(&codes, "codes", nil, "Codes (eg. us-east-1)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
