package cli

import (
	"fmt"
	"strings"
	"sync"

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
