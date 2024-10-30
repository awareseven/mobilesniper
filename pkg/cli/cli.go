package cli

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var maxConcurrency int // maximum number of concurrent Go-routines
var noColor bool       // don't use ANSI colors
var verbose bool       // verbose mode

var rootCmd = &cobra.Command{
	Use:   "mobilesniper",
	Short: "A pentesting tool for 5G mobile networks.",
	Long:  "MobileSniper is a CLI application for performing various pentesting tasks specialicied on 5G mobile networks.",
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewProgressBar(count int, desc string) (*progressbar.ProgressBar, *ProgressLogger) {

	var theme progressbar.Theme
	var description string

	if noColor {
		theme = progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}
		description = desc
	} else {
		theme = progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "][green]",
		}
		description = fmt.Sprintf("[red]%s[green]", desc)
	}

	bar := progressbar.NewOptions(count,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(25),
		progressbar.OptionSetElapsedTime(true),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSpinnerType(21),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
		progressbar.OptionSetTheme(theme),
	)

	// Create the custom logger and set it as the default logger
	progressLogger := NewProgressLogger(bar)
	log.SetOutput(progressLogger)

	// Periodically refresh the progress bar to update the elapsed time
	go func() {
		for {
			time.Sleep(time.Second)
			bar.RenderBlank() // Refresh the progress bar without changing progress
		}
	}()

	return bar, progressLogger
}

func init() {
	rootCmd.AddCommand(enumCmd)
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(scanCmd)

	rootCmd.PersistentFlags().IntVarP(
		&maxConcurrency, "max-goroutines", "c", 128, "Maximum number of concurrent Go-routines",
	)
	rootCmd.PersistentFlags().BoolVar(
		&noColor, "no-color", false, "Don't use ANSI colors",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false, "Verbose mode",
	)
}
