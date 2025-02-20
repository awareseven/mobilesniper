package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/cli/analyze"
	"github.com/awareseven/mobilesniper/cli/core"
	"github.com/awareseven/mobilesniper/cli/enum"
	"github.com/awareseven/mobilesniper/cli/scan"
)

var rootCmd = &cobra.Command{
	Use:   "mobilesniper",
	Short: "A pentesting tool for 5G mobile networks.",
	Long:  "MobileSniper is a CLI application for performing various pentesting tasks specialicied on 5G mobile networks.",
}

func init() {
	rootCmd.AddCommand(enum.EnumCmd)
	rootCmd.AddCommand(analyze.AnalyzeCmd)
	rootCmd.AddCommand(scan.ScanCmd)

	rootCmd.PersistentFlags().IntVarP(
		&core.MaxConcurrency, "max-goroutines", "c", 128, "Maximum number of concurrent Go-routines",
	)
	rootCmd.PersistentFlags().BoolVar(
		&core.NoColor, "no-color", false, "Don't use ANSI colors",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&core.Verbose, "verbose", "v", false, "Verbose mode",
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
