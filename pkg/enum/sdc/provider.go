package enum

import (
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/ipv4"

	models "github.com/awareseven/mobilesniper/pkg/models/soap"
	"github.com/awareseven/mobilesniper/pkg/utils"
)

func DiscoverSDCDevices(interfaceName string, sdcChan chan<- models.ReceiveSOAPMessage, wg *sync.WaitGroup, verbose bool) error {
	defer wg.Done()

	// Find correct source interface
	srcInterface, srcIP, err := utils.GetNetInterfaceAndIPNet(interfaceName, verbose)
	if err != nil {
		log.Fatalf("%v", err)
	}
	utils.LogVerbosef(verbose, "Will bind to %s on interface %s", srcIP.IP, srcInterface.Name)

	// Create UDP addresses to send and recive discovery data
	srcAddr := &net.UDPAddr{IP: srcIP.IP, Port: 4747}
	discoveryAddr := &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 3702}

	multicastReceiver, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		log.Fatalf("Failed to create UDP connection: %v", err)
	}
	defer multicastReceiver.Close()

	// Join the multicast group
	udpConn := ipv4.NewPacketConn(multicastReceiver)
	if err := udpConn.JoinGroup(srcInterface, discoveryAddr); err != nil {
		log.Fatalln(err)
	}
	defer udpConn.LeaveGroup(srcInterface, discoveryAddr)

	// Create SOAP message - Probe
	probe := models.NewSOAPMessage(
		models.NewProbeSOAPHeader(),
		models.NewProbeSOAPBody(),
	)

	// Marshal the Probe message to XML
	xmlBytes, err := probe.XMLMarshal()
	if err != nil {
		log.Fatalln(err)
	}

	// Send the Probe message
	_, err = udpConn.WriteTo(xmlBytes, nil, discoveryAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Send Probe on %s", discoveryAddr)

	// Set deadline for connection to five second
	if err := udpConn.SetReadDeadline(time.Now().Add(time.Second * 5)); err != nil {
		log.Fatal(err)
	}

	// Read UDP multicast response as long as there is a communication channel open
	for {
		// a buffer of 8Mb should be large enought to store probe matches
		buffer := make([]byte, 8192)
		msgLength, _, srcAddr, err := udpConn.ReadFrom(buffer) // ignore IMCP control message and source address
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				// No UDP responses to read will result into an i/o timeout
				err = nil // this isn't a error, it just indicates that all responses are already read
				break
			}
			log.Fatal(err)
		}
		utils.LogVerbosef(verbose, "Received %d bytes anwser from %s:\n%v", msgLength, srcAddr, string(buffer))
		res, err := models.XMLUnmarshal[models.ProbeMatchBody](buffer)
		body, ok := res.SOAPBody.Payload.(models.ProbeMatchBody)

		if err == nil && ok {
			sdcChan <- res
			// this is the HTTP server that expects SDC SOAP messages
			log.Printf("Found SDC endpoint: %s", body.GetXAddrs())
		} else {
			log.Printf(
				"Received %d bytes anwser from %s that could not interpreted as a 'ProbeMatch'.", msgLength, srcAddr,
			)
			utils.LogVerbosef(verbose, "XML unmarshal of %d bytes anwser from %s error: %v", msgLength, srcAddr, err)
		}
	}

	return err
}
