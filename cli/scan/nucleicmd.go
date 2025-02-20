package scan

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/cli/core"
	"github.com/awareseven/mobilesniper/pkg/scan"
	utils "github.com/awareseven/mobilesniper/pkg/utils"
)

var nucleiCmd = &cobra.Command{
	Use:   "nuclei <network range or single IP>",
	Short: "Run a nuclei scan on the given target.",
	Long:  "This command runs nuclei to perform a vulnerability scan.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cidrOrIP := args[0]
		_, err := utils.GetIPsInCIDR(cidrOrIP)
		if err != nil {
			panic(err)
		}

		bar, _ := core.NewProgressBar(1, fmt.Sprintf("Scanning: %s", cidrOrIP))
		defer bar.Finish()

		var wg sync.WaitGroup
		resultsChan := make(chan output.ResultEvent)

		wg.Add(1)
		go scan.RunNucleiScan(cidrOrIP, resultsChan, &wg, core.Verbose)

		go func() {
			wg.Wait()
			close(resultsChan)
		}()

		for event := range resultsChan {
			bar.ChangeMax(bar.GetMax() + 1)
			log.Printf("%s:%s - %v", event.IP, event.Port, event.ExtractedResults)

			time.Sleep(100 * time.Millisecond)
			bar.Add(1)
		}
	},
}
