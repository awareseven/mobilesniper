package analyze

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/cli/core"
	"github.com/awareseven/mobilesniper/pkg/analyze"
)

var pcapCmd = &cobra.Command{
	Use:   "pcap",
	Short: "Analyze .pcap-file",
	Long:  "Analyze Wireshark traffic capturing (.pcap file) related to 5G communication.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		pcapFile := args[0]
		openapiPath, _ := cmd.Flags().GetString("openapi")

		bar, _ := core.NewProgressBar(1, fmt.Sprintf("Analyzing: %s", pcapFile))
		defer bar.Finish()

		if _, err := os.Stat(pcapFile); os.IsNotExist(err) {
			log.Printf("File doesn't exist: %s", pcapFile)
			return
		}

		nfr, err := analyze.AnalyzePcap(pcapFile, openapiPath, core.Verbose)
		if err != nil {
			log.Println(err)
		}

		if len(nfr) != 0 {
			bar.ChangeMax(len(nfr))
		} else {
			log.Println("Could not find any netfowk function by traffic analysis.")
		}

		for _, result := range nfr {
			log.Println(result.String())

			time.Sleep(100 * time.Millisecond)
			bar.Add(1)
		}
	},
}
