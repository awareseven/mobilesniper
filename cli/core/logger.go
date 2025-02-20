package core

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type ProgressLogger struct {
	bar  *progressbar.ProgressBar
	lock sync.Mutex
}

func NewProgressLogger(bar *progressbar.ProgressBar) *ProgressLogger {
	return &ProgressLogger{
		bar: bar,
	}
}

func (pl *ProgressLogger) Write(p []byte) (n int, err error) {
	pl.lock.Lock()
	defer pl.lock.Unlock()

	// Split the log message into lines
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			// Clear the current progress bar from the terminal
			pl.bar.Clear()

			// Print the log message line
			fmt.Println(line)

			// Re-render the progress bar after each line
			pl.bar.RenderBlank()
		}
	}

	return len(p), nil
}

func NewProgressBar(count int, desc string) (*progressbar.ProgressBar, *ProgressLogger) {

	var theme progressbar.Theme
	var description string

	if NoColor {
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
