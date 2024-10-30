package cli

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/pkg/enum"
	"github.com/awareseven/mobilesniper/pkg/models"
	utils "github.com/awareseven/mobilesniper/pkg/utils"
)

var enumCmd = &cobra.Command{
	Use:   "enum",
	Short: "Enumeration commands",
	Long:  `This command group contains commands related to enumeration.`,
}

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

		bar, _ := NewProgressBar(1, fmt.Sprintf("Discover Services: %s", cidrOrIP))

		var wg sync.WaitGroup
		targetChan := make(chan models.Target)

		wg.Add(1)
		hostTimeout, _ := cmd.Flags().GetString("host-timeout")
		go enum.DiscoverOpenPorts(cidrOrIP, targetChan, &wg, maxConcurrency, hostTimeout, verbose, "-sV")

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

		bar, _ := NewProgressBar(1, fmt.Sprintf("Discover Network Functions: %s", cidrOrIP))

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
			go enum.DiscoverOpenPorts(cidrOrIP, targetChan, &targetWg, maxConcurrency, hostTimeout, verbose)

			go func() {
				targetWg.Wait()
				close(targetChan)
			}()

		}

		for target := range targetChan {
			bar.ChangeMax(bar.GetMax() + len(files))

			nfrWg.Add(1)
			go enum.DiscoverNetworkFunctions(target, openapiPath, nfrChan, &nfrWg, maxConcurrency, verbose)
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

func init() {
	nfsCmd.Flags().String(
		"openapi", "assets/5GC-APIs", "Path to 3GPP OpenAPI definitions of 5G network functions",
	)
	nfsCmd.Flags().Float64P(
		"threshold", "t", 70.0, "The threshold of accurancy a NF should be considered as detected.",
	)

	enumCmd.PersistentFlags().String(
		"host-timeout", "4m", "The time before give up a scan on a single host.",
	)

	enumCmd.AddCommand(servicesCmd)
	enumCmd.AddCommand(nfsCmd)
}
