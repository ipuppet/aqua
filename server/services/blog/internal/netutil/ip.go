package netutil

import (
	"net"
)

// ParseIP parses a string IP address and returns its net.IP form and protocol version (4 or 6).
func ParseIP(s string) (net.IP, int) {
	parsedIP := net.ParseIP(s)
	if parsedIP == nil {
		return nil, 0
	}

	ipv4 := parsedIP.To4()
	if ipv4 != nil {
		return ipv4, 4
	}
	ipv6 := parsedIP.To16()
	if ipv6 != nil {
		return ipv6, 6
	}

	return nil, 0
}
