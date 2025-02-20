package enum

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/cli/core"
	"github.com/awareseven/mobilesniper/pkg/enum"
	"github.com/awareseven/mobilesniper/pkg/models"
	utils "github.com/awareseven/mobilesniper/pkg/utils"
)

var nfsCmd = &cobra.Command{
	Use:   "nf <network range>",
	Short: "Probe network to identify specific network function",
	Long:  "This command enumerates 3GPP OpenAPI definition to identifies 5G network function within the given network range or IP address.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cidrOrIP := args[0]

		openapiPath, _ := cmd.Flags().GetString("openapi")
		threshold, _ := cmd.Flags().GetFloat64("threshold")
		files, _ := os.ReadDir(openapiPath)

		bar, _ := core.NewProgressBar(1, fmt.Sprintf("Discover Network Functions: %s", cidrOrIP))

		var targetWg, nfrWg sync.WaitGroup
		targetChan := make(chan models.Target)
		nfrChan := make(chan models.NetworkFunctionResult)

		mayIpAndPort := strings.Split(cidrOrIP, ":")
		if len(mayIpAndPort) == 2 {

			port, _ := strconv.Atoi(mayIpAndPort[1])
			target := models.Target{
				IP:       mayIpAndPort[0],
				Port:     port,
				Protocol: "tcp",
			}

			go func() {
				targetChan <- target
				close(targetChan)
			}()

		} else {
			_, err := utils.GetIPsInCIDR(cidrOrIP)
			if err != nil {
				panic(err)
			}

			targetWg.Add(1)
			hostTimeout, _ := cmd.Flags().GetString("host-timeout")
			go enum.DiscoverOpenPorts(cidrOrIP, targetChan, &targetWg, core.MaxConcurrency, hostTimeout, core.Verbose)

			go func() {
				targetWg.Wait()
				close(targetChan)
			}()

		}

		for target := range targetChan {
			bar.ChangeMax(bar.GetMax() + len(files))

			nfrWg.Add(1)
			go enum.DiscoverNetworkFunctions(target, openapiPath, nfrChan, &nfrWg, core.MaxConcurrency, core.Verbose)
		}

		go func() {
			nfrWg.Wait()
			close(nfrChan) // Close the nfr channel after all NF discovery is done
			bar.Finish()   // Ensure progress bar finishes after all operations
		}()

		for nfr := range nfrChan {
			// Log only network function that have a detection rate over 50%
			// and under exactly 100%. A accurancy of 100% is most likly a false positive.
			//
			// E.g. a python http server was always detected as some NFs
			//
			if nfr.Accuracy > threshold && nfr.Accuracy < 100.00000 {
				log.Println(nfr.String())
			}

			time.Sleep(100 * time.Millisecond)
			bar.Add(1)
		}
	},
}
