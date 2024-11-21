package server

import (
	"fmt"
	"net"
)

func NormalizeIPv4(address string) (string, error) {
	ip := net.ParseIP(address)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address: %s", address)
	}

	// Check if the IP is an IPv6 address mapped to IPv4
	if ip.To4() != nil {
		return ip.String(), nil
	}

	// Handle special case for localhost
	if ip.IsLoopback() {
		return "127.0.0.1", nil
	}

	return "", fmt.Errorf("not an IPv4-mapped IPv6 address: %s", address)
}
