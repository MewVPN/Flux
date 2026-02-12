package util

import (
	"net"
)

// DetectPublicIP tries to find a non-loopback IPv4 address
func DetectPublicIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
	}

	for _, iface := range ifaces {
		// interface must be up
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not IPv4
			}

			if ip.IsLoopback() {
				continue
			}

			// found a valid IPv4
			return ip.String()
		}
	}

	// fallback
	return "127.0.0.1"
}
