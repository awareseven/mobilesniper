package enum

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/awareseven/mobilesniper/cli/core"
	enum "github.com/awareseven/mobilesniper/pkg/enum/sdc"
	"github.com/awareseven/mobilesniper/pkg/models"
	soap "github.com/awareseven/mobilesniper/pkg/models/soap"
	"github.com/awareseven/mobilesniper/pkg/utils"
)

var sdcCmd = &cobra.Command{
	Use:   "sdc",
	Short: "SDC / ISO 11703 related commands",
	Long:  `This command group contains commands related to SDC / ISO 11703.`,
}

var providerCmd = &cobra.Command{
	Use:   "provider <network interface>",
	Short: "Enumerate SDC devices.",
	Long:  `This command performs a SDC discovery process via UDP multicast to discover medical devices.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		ifi := args[0]
		bar, _ := core.NewProgressBar(1, fmt.Sprintf("Discover SDC Devices: %s", ifi))
		defer bar.Finish()

		var wg sync.WaitGroup
		sdcChan := make(chan soap.ReceiveSOAPMessage)

		wg.Add(1)
		go enum.DiscoverSDCDevices(ifi, sdcChan, &wg, core.Verbose)

		go func() {
			wg.Wait()
			close(sdcChan) // Ensure channel closure after all operations are complete
		}()

		for soapMsg := range sdcChan {
			bar.ChangeMax(bar.GetMax() + 1)

			probeMatch, ok := soapMsg.SOAPBody.Payload.(soap.ProbeMatchBody)
			if !ok {
				log.Println("Received SOAP message thas wan't a 'ProbeMatch'.")
				utils.LogVerbosef(core.Verbose, "Received:\n%v", soapMsg)
				continue
			}
			bar.Add(1)

			// call HTTP server with a "Get" message that should return device information
			soapGetMsg := soap.NewSOAPMessage(
				soap.NewGetSOAPHeader(probeMatch.GetXAddrs()),
				soap.NewGetSOAPBody(),
			)

			// Marshal the Probe message to XML
			xmlBytes, err := soapGetMsg.XMLMarshal()
			if err != nil {
				log.Println(err)
				continue // continue with next received probe match
			}

			// actually perform the HTTP request with given porbe match address
			_, bodyBytes, err := utils.DoSdcXClientSOAPPost(probeMatch.GetXAddrs(), xmlBytes, core.Verbose)
			if err != nil {
				log.Println(err) // on mTLS enforced devices, this will return a mTLS mismatch error
				continue         // continue with next received probe match
			}

			// try to parse the returned SOAP message
			getResp, _ := soap.XMLUnmarshal[soap.GetResponseBody](bodyBytes)

			// TODO: SDC communication to perform (more) malicious actions
			// -> successfully connection indicates that the device doesn't enfore mTLS
			// -> without mTLS there is no other authentication method (as soon as we know)
			device, _ := models.CreateSDCDevicebyGetResponse(getResp)
			log.Printf("Endpoint states it is %s", device.String())

			time.Sleep(100 * time.Millisecond)
			bar.Add(1)
		}
	},
}

var consumerCmd = &cobra.Command{
	Use:   "consumer <network interface>",
	Short: "Enumerate SDC consumers.",
	Long:  `This command performs a SDC discovery process via UDP multicast to discover SDC consumers.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var serverCrt, serverKey, caCrt *string

		ifi := args[0]
		crtFlag, _ := cmd.Flags().GetString("tlsCrt")
		keyFlag, _ := cmd.Flags().GetString("tlsKey")
		caFlag, _ := cmd.Flags().GetString("tlsCA")

		if crtFlag != "" && keyFlag != "" && caFlag != "" {
			serverCrt = &crtFlag
			serverKey = &keyFlag
			caCrt = &caFlag
		} else {
			log.Println(
				"Continue without TLS. If you intended to use TLS than provide `tlsCrt`, `tlsKey` and `tlsCA` flag.",
			)
		}

		bar, _ := core.NewProgressBar(1, fmt.Sprintf("Discover SDC Consumers: %s", ifi))
		defer bar.Finish()

		var wg sync.WaitGroup
		// sdcChan := make(chan soap.ProbeMatch)

		wg.Add(1)
		go enum.DiscoverSDCConsumer(ifi, &wg, caCrt, serverCrt, serverKey, core.Verbose)
		wg.Wait()
	},
}
