package enum

import (
	"crypto/tls"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/ipv4"

	models "github.com/awareseven/mobilesniper/pkg/models/soap"
	"github.com/awareseven/mobilesniper/pkg/utils"
)

func DiscoverSDCConsumer(interfaceName string, wg *sync.WaitGroup, caCrt, serverCrt, serverKey *string, verbose bool) error {
	defer wg.Done()

	useTLS := false
	if serverCrt != nil && serverKey != nil && caCrt != nil {
		useTLS = true
	}

	// Find correct source interface
	srcInterface, srcIP, err := utils.GetNetInterfaceAndIPNet(interfaceName, verbose)
	if err != nil {
		log.Fatalln(err)
	}
	utils.LogVerbosef(verbose, "Will bind to %s on interface %s", srcIP.IP, srcInterface.Name)

	//
	// Step 1: Send Provider Hello and receive Probe message via UDP multicast
	//

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

	// Set deadline for connection to five second
	if err := udpConn.SetReadDeadline(time.Now().Add(time.Second * 5)); err != nil {
		log.Fatal(err)
	}

	// Listen for incoming connections
	var tcpList net.Listener
	servAddr := &net.TCPAddr{IP: srcIP.IP, Port: 4748}

	if useTLS {
		// Start TCP server with TLS encryption
		cert, err := tls.LoadX509KeyPair(*serverCrt, *serverKey)
		if err != nil {
			log.Fatalf("could not load TLS config: %v", err)
		}

		tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

		tcpList, err = tls.Listen("tcp", servAddr.AddrPort().String(), tlsCfg)
		if err != nil {
			log.Fatal(err)
		}
		defer tcpList.Close()
	} else {
		// Start TCP server without encryption
		tcpList, err = net.ListenTCP("tcp", servAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer tcpList.Close()
	}

	// Create SOAP message - Hello
	hello := models.NewSOAPMessage(
		models.NewHelloSOAPHeader(),
		models.NewHelloSOAPBody(servAddr),
	)

	// Marshal the Probe message to XML
	xmlBytes, err := hello.XMLMarshal()
	if err != nil {
		log.Fatalln(err)
	}

	// Send the Hello message
	_, err = udpConn.WriteTo(xmlBytes, nil, discoveryAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Send Hello on %s", discoveryAddr)

	//
	// ----------- This error always appears when dealing with the sdcX Consumer -----------------------
	//
	//	- The consumer receives the hello mesage, but discards the request.
	//	- The `ConsumerCLI -debug -tls` log confirm that the consumer does not send any response.
	//	- Another approach could be to listen for probe messages. But this would not allow an active
	//	  dicovery procedure of SDC consumers.
	//
	//	- ConpleteConsumer.Log of `ConsumerCLI -debug -tls`:
	//
	// [7][2025-02-17 11:59:37:097] UDPResponseMessageConsumer: <?xml version="1.0" encoding="UTF-8"?> ...
	// [6][2025-02-17 11:59:37:098] RoutingLinkingLayer: Routing message 1172
	// [6][2025-02-17 11:59:37:098] SchemaValidationLinkingLayer: INFORMATIONAL: Received network transport with ID 1172
	// [6][2025-02-17 11:59:37:099] SOAPProcessingLinkingLayer:
	// 								Received incoming soap envelope transport message with ID 1173
	// [6][2025-02-17 11:59:37:100] ConsumerDiscoveryServiceLinkingLayer: Received hello message with ID 1175
	// Received hello from unknown device: urn:uuid:08960b15-a692-462c-8cb7-527c3a213e86
	//
	// -------------------------------------------------------------------------------------------------
	//

	// Read UDP multicast response as long as there is a communication channel open
	for {
		// a buffer of 8Mb should be large enought to store hello responses
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
		res, err := models.XMLUnmarshal[models.ProbeBody](buffer)
		body, ok := res.SOAPBody.Payload.(models.ProbeBody)

		if err == nil && ok {
			// this is the HTTP server that expects SDC SOAP messages
			log.Println(body)
		} else {
			log.Printf(
				"Received %d bytes anwser from %s that could not interpreted as a 'Probe'.", msgLength, srcAddr,
			)
			utils.LogVerbosef(verbose, "XML unmarshal of %d bytes anwser from %s error: %v", msgLength, srcAddr, err)
		}
	}

	//
	// Step 2: Receive Get message and anwser via HTTP SOAP API
	//

	//! Caution: Will never terminate!!
	//	TODO: Implement break condition: timeout? fixed time?
	for {
		// Accept incoming connections
		conn, err := tcpList.Accept()
		if err != nil {
			log.Println(err)
		}

		// Handle client connection in a goroutine
		wg.Add(1)
		go handleClient(conn, wg)
	}

	//
	// Step 3: Receive all Subscribe message and anwser via HTTP SOAP API
	//

	// TODO: Implement Step 3

	//
	// Step 4: Receive GetMdib message and anwser via HTTP SOAP API
	//

	// TODO: Implement Step 4

	//
	// Step 5: Send SubscriptionResponse messages to consumer via HTTP SOAP API
	//

	// TODO: Implement Step 5

	return err
}

func handleClient(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()

	// Create a buffer to read data into
	buffer := make([]byte, 1024)

	for {
		// Read data from the client
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error:", err)
			return
		}

		// Process and use the data (here, we'll just print it)
		log.Printf("Received: %s\n", buffer[:n])
	}
}
