package utils

import (
	"fmt"
	"log"
	"net"
)

func GetIPsInCIDR(cidrOrIP string) ([]string, error) {
	ip := net.ParseIP(cidrOrIP)
	if ip != nil {
		return []string{cidrOrIP}, nil
	}

	ip, ipNet, err := net.ParseCIDR(cidrOrIP)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR notation: %v", err)
	}

	var ips []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); nextIP(ip) {
		// Check if IP is the network or broadcast address
		if !ipNet.Contains(ip) || ip.Equal(ipNet.IP) || isBroadcast(ip, ipNet) {
			continue
		}
		ips = append(ips, ip.String())
	}

	return ips, nil
}

func nextIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
func isBroadcast(ip net.IP, ipNet *net.IPNet) bool {
	broadcast := make(net.IP, len(ip))
	for i := range ip {
		broadcast[i] = ip[i] | ^ipNet.Mask[i]
	}
	return ip.Equal(broadcast)
}

func GetNetInterfaceAndIPNet(interfaceName string, verbose bool) (*net.Interface, *net.IPNet, error) {
	var srcInterface *net.Interface
	var srcIP *net.IPNet

	// Get source interface either by given name or
	// the first none-loopback device.

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Failed to enumerate network interfaces: %v", err)
	}
	LogVerbosef(verbose, "Network interfaces: %v", interfaces)

	// iterate over available network interfaces to match the given interface argument
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Printf("Failed to enumerate network addresses: %v", err)
			continue
		}
		LogVerbosef(verbose, "Interface %s addresses: %v", i.Name, addrs)

		// just use the first address of the network
		srcIP = addrs[0].(*net.IPNet)
		LogVerbosef(verbose, "Got source IP %v on %s.", srcIP.IP, i.Name)

		if interfaces != nil {
			if i.Name == interfaceName {
				srcInterface = &i
				break
			}
		} else {
			if !srcIP.IP.IsLoopback() {
				srcInterface = &i
				break
			}
		}
	}

	if srcIP.IP == nil || srcInterface == nil {
		log.Fatalln("Failed to detect none-loopback interface to bind.")
	}

	return srcInterface, srcIP, err
}
