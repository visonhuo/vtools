package cmd

import (
	"math"
	"strings"
)

func isIPv6Protocol(protocol string) bool {
	return strings.HasSuffix(protocol, "6")
}

func isValidPort(port int) bool {
	if port < 0 || port > math.MaxUint16 {
		return false
	}
	return true
}
