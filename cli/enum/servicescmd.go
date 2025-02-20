package enum

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/cli/core"
	"github.com/awareseven/mobilesniper/pkg/enum"
	"github.com/awareseven/mobilesniper/pkg/models"
	utils "github.com/awareseven/mobilesniper/pkg/utils"
)

var servicesCmd = &cobra.Command{
	Use:   "services <network range or single IP>",
	Short: "Perform a port scan on a given network",
	Long:  `This command performs a port scan with service discovery on a given network range by nmap`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cidrOrIP := args[0]
		_, err := utils.GetIPsInCIDR(cidrOrIP)
		if err != nil {
			panic(err)
		}

		bar, _ := core.NewProgressBar(1, fmt.Sprintf("Discover Services: %s", cidrOrIP))

		var wg sync.WaitGroup
		targetChan := make(chan models.Target)

		wg.Add(1)
		hostTimeout, _ := cmd.Flags().GetString("host-timeout")
		go enum.DiscoverOpenPorts(cidrOrIP, targetChan, &wg, core.MaxConcurrency, hostTimeout, core.Verbose, "-sV")

		go func() {
			wg.Wait()
			close(targetChan) // Ensure channel closure after all operations are complete
			bar.Finish()      // Ensure progress bar finishes after all operations
		}()

		for target := range targetChan {
			bar.ChangeMax(bar.GetMax() + 1)
			log.Println(target)

			time.Sleep(100 * time.Millisecond)
			bar.Add(1)
		}
	},
}
