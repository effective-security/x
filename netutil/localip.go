package netutil

import (
	"net"
	"time"

	"github.com/cockroachdb/errors"
)

// IsPrivateAddress reports whether address is loopback, link-local, or RFC1918/ULA.
//
// Deprecated: use net.ParseIP(address) with net.IP.IsPrivate, IsLoopback,
// and IsLinkLocalUnicast. IsPrivate does not treat loopback or link-local
// as private; this helper does.
func IsPrivateAddress(address string) (bool, error) {
	ipAddress := net.ParseIP(address)
	if ipAddress == nil {
		return false, errors.New("address is not valid")
	}

	return ipAddress.IsPrivate() || ipAddress.IsLoopback() || ipAddress.IsLinkLocalUnicast(), nil
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.WithStack(err)
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("unable to resolve local IP address")
}

// WaitForNetwork will wait until the local IP is available or timeout ocurred
func WaitForNetwork(d time.Duration) (ipaddr string, err error) {
	ipaddr, err = GetLocalIP()
	if err != nil {
		cutoff := time.Now().Add(d)
		for cutoff.After(time.Now()) {
			ipaddr, err = GetLocalIP()
			if err == nil {
				break
			}
			time.Sleep(time.Second)
		}
	}
	return
}
