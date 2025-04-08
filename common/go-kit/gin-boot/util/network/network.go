package network

import (
	"errors"
	"net"
	"strings"
)

// IsVirtualInterface return true if interface is virtual
func IsVirtualInterface(iface string) bool {
	return strings.Contains(iface, ":")
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP(interfaceName string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {

		if interfaceName != "" && iface.Name != interfaceName {
			continue
		}

		if (iface.Flags&net.FlagLoopback == 0) &&
			(iface.Flags&net.FlagUp != 0) &&
			!IsVirtualInterface(iface.Name) {

			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}
			for _, address := range addrs {
				// check the address type and if it is not a loopback the display it
				if ipnet, ok := address.(*net.IPNet); ok {
					if ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() {
						return ipnet.IP.String(), nil
					}
				}
			}
		}
	}

	return "", errors.New("machine have no available network interface OR named interface have no ip address")
}
