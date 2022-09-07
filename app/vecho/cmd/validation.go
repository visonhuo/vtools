package cmd

import (
	"math"
	"strings"
)

func isValidProtocol(protocol string) bool {
	switch protocol {
	case "tcp", "udp", "tcp4", "udp4", "tcp6", "udp6":
		return true
	default:
		return false
	}
}

func isIPv6Protocol(protocol string) bool {
	return strings.HasSuffix(protocol, "6")
}

func isValidPort(port int) bool {
	if port < 0 || port > math.MaxUint16 {
		return false
	}
	return true
}
