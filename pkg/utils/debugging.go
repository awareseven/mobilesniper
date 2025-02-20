package utils

import (
	"fmt"
	"log"
)

// Helper function to print logs on debugging.
func LogVerbosef(verbose bool, s string, v ...any) {
	if verbose {
		format := fmt.Sprintf(s, v...)
		log.Println(format)
	}
}
