package analyze

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	"github.com/awareseven/mobilesniper/pkg/models"
)

func AnalyzePcap(pcapFile, openapiPath string, verbose bool) ([]models.NetworkFunctionResult, error) {
	handle, err := pcap.OpenOffline(pcapFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open PCAP file: %v", err)
	}
	defer handle.Close()

	var retVal []models.NetworkFunctionResult

	// Create a packet source to read packets from the PCAP file.
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {

		// Skip non-IP traffic
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}

		ipPacket, _ := ipLayer.(*layers.IPv4)
		if ipPacket.Protocol != layers.IPProtocolTCP {
			continue // Skip non-TCP packets
		}

		// Extract the TCP layer
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		tcpPacket, _ := tcpLayer.(*layers.TCP)

		dstPort := int(tcpPacket.DstPort)
		dstIP := ipPacket.DstIP
		target := models.Target{IP: dstIP.String(), Port: dstPort, Protocol: "tcp"}

		// Extract the application layer (payload) from the packet.
		applicationLayer := packet.ApplicationLayer()
		if applicationLayer == nil {
			continue
		}
		payload := applicationLayer.Payload()

		// Check if the payload contains HTTP data.
		if !isHTTPPayload(payload) {
			continue
		}

		// Extract HTTP request information.
		reqMethod, reqPath := extractHTTPRequestInfo(payload)
		if reqMethod == "" || reqPath == "" {
			continue
		}

		// Enumerare OpenAPI definitions to detect network functions
		err = filepath.Walk(openapiPath, func(path string, info os.FileInfo, err error) error {
			// Validate OpenAPI definitions
			openapi, err := models.ValidateOpenAPIFile(path)
			if err != nil {
				return nil // returning error leads to aporting the Walk function
			}

			// Iterate over all paths defined in the OpenAPI specification
			for apiPath, apiMethods := range openapi.Paths {
				// Prefix path with server root if server root exists
				apiPathWithRoot := apiPath
				if len(openapi.Servers) != 0 {
					apiPathWithRoot = openapi.Servers[0].URL + apiPath
				}

				// Iterate over all HTTP methods (GET, POST, etc.) for each path
				for apiMethod, _ := range apiMethods {
					if matchReq(reqPath, apiPathWithRoot, reqMethod, apiMethod) {
						result := models.NetworkFunctionResult{
							Accuracy:        100.0,
							NetworkFunction: openapi.Info.Title,
							Target:          target,
						}

						if !models.ContainsNFResult(retVal, result) {
							retVal = append(retVal, result)

							if verbose {
								log.Printf("%s - %s: %s", result.String(), reqMethod, reqPath)
							}
						}

						return nil // continue with next NF
					}
				}
			}

			return err
		})

	}

	return retVal, err
}

// matchReq checks if the request match the definied API request
func matchReq(reqPath, apiPath, reqMethod, apiMethod string) bool {
	// check if method is the same
	if strings.ToUpper(reqMethod) != strings.ToUpper(apiMethod) {
		return false
	}

	// check if path lenghts match
	reqSegments := strings.Split(strings.Trim(reqPath, "/"), "/")
	apiSegments := strings.Split(strings.Trim(apiPath, "/"), "/")

	if len(reqSegments) != len(apiSegments) {
		if len(reqSegments) == len(apiSegments)-1 {
			// apiRoot is none / server root path
			apiSegments = apiSegments[1:]
		} else {
			return false
		}
	}

	// check segments
	for i, regSeg := range reqSegments {
		apiSeg := apiSegments[i]
		if strings.HasPrefix(apiSeg, "{") && strings.HasSuffix(apiSeg, "}") {
			// This segment is a path parameter; accept any value.
			continue
		}

		// path is different
		if regSeg != apiSeg {
			return false
		}
	}

	return true
}

// isHTTPPayload checks if the payload appears to be an HTTP request or response.
func isHTTPPayload(payload []byte) bool {
	httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "CONNECT", "TRACE"}

	payloadStr := string(payload)
	for _, method := range httpMethods {
		if strings.HasPrefix(payloadStr, method+" ") {
			return true
		}
	}
	return false
}

// extractHTTPRequestInfo extracts the HTTP method and path from the payload.
func extractHTTPRequestInfo(payload []byte) (method string, path string) {
	lines := strings.Split(string(payload), "\n")
	if len(lines) == 0 {
		return "", ""
	}

	// The first line should contain the request line: METHOD PATH PROTOCOL
	requestLineParts := strings.Split(strings.TrimSpace(lines[0]), " ")
	if len(requestLineParts) < 2 {
		return "", ""
	}

	return requestLineParts[0], requestLineParts[1]
}
