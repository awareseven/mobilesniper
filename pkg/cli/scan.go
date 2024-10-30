package cli

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/pkg/models"
	"github.com/awareseven/mobilesniper/pkg/scan"
	utils "github.com/awareseven/mobilesniper/pkg/utils"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Perfrom a vulnerability scan.",
	Long:  "This command hroup contains commands to discovery network vulnerabilities.",
}

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

		bar, _ := NewProgressBar(1, fmt.Sprintf("Scanning: %s", cidrOrIP))
		defer bar.Finish()

		var wg sync.WaitGroup
		resultsChan := make(chan output.ResultEvent)

		wg.Add(1)
		go scan.RunNucleiScan(cidrOrIP, resultsChan, &wg, verbose)

		go func() {
			wg.Wait()
			close(resultsChan)
		}()

		for event := range resultsChan {
			bar.ChangeMax(bar.GetMax() + 1)
			log.Printf("%s:%d - %v", event.IP, event.Port, event.ExtractedResults)

			time.Sleep(100 * time.Millisecond)
			bar.Add(1)
		}
	},
}

var nessusCmd = &cobra.Command{
	Use:   "nessus <network range or single IP>",
	Short: "Run a nessus scan on the given target.",
	Long:  "This command runs nessus to perform a vulnerability scan.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cidrOrIP := args[0]
		_, err := utils.GetIPsInCIDR(cidrOrIP)
		if err != nil {
			panic(err)
		}

		nessus_username, err := cmd.Flags().GetString("username")
		nessus_password, err := cmd.Flags().GetString("password")
		nessus_url, err := cmd.Flags().GetString("url")

		if err != nil {
			panic(err)
		}

		config := models.NessusConf{
			Username: nessus_username,
			Password: nessus_password,
			URL:      nessus_url,
		}

		bar, _ := NewProgressBar(1, fmt.Sprintf("Scanning: %s", cidrOrIP))
		defer bar.Finish()

		var wg sync.WaitGroup

		wg.Add(1)
		go scan.RunNessusScan(cidrOrIP, config, &wg, verbose)
		wg.Wait()
	},
}

func init() {
	nessusCmd.PersistentFlags().String(
		"username", "", "Nessus username to authenticate",
	)
	nessusCmd.PersistentFlags().String(
		"password", "", "Nessus password for given user to authenticate.",
	)
	nessusCmd.PersistentFlags().String(
		"url", "https://127.0.0.1:8834", "URL to Nessus Web-UI",
	)

	scanCmd.AddCommand(nucleiCmd)
	scanCmd.AddCommand(nessusCmd)
}
