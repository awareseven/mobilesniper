package enum

import (
	"fmt"
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
