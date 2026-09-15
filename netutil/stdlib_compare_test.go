package netutil_test

import (
	"net"
	"testing"

	"github.com/effective-security/x/netutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const benchmarkIPAddress = "192.168.100.200"

var benchmarkPrivateResult bool

func TestIsPrivateAddressStdlibCompatibility(t *testing.T) {
	addresses := []string{
		"127.0.0.2",
		"10.0.0.1",
		"172.16.0.1",
		"192.168.0.1",
		"169.254.1.1",
		"::1",
		"fc00::1",
		"fe80::1",
		"8.8.8.8",
		"2001:4860:4860::8888",
	}
	for _, address := range addresses {
		t.Run(address, func(t *testing.T) {
			actual, err := netutil.IsPrivateAddress(address)
			require.NoError(t, err)
			assert.Equal(t, isPrivateStdlib(net.ParseIP(address)), actual)
		})
	}
}

func BenchmarkIsPrivateAddressVsStdlib(b *testing.B) {
	ip := net.ParseIP(benchmarkIPAddress)
	cidrs := legacyPrivateCIDRs()
	b.Run("LegacyCIDRScan", func(b *testing.B) {
		for b.Loop() {
			benchmarkPrivateResult = isPrivateLegacy(benchmarkIPAddress, cidrs)
		}
	})
	b.Run("CurrentWrapper", func(b *testing.B) {
		for b.Loop() {
			benchmarkPrivateResult, _ = netutil.IsPrivateAddress(benchmarkIPAddress)
		}
	})
	b.Run("StdlibParsedIP", func(b *testing.B) {
		for b.Loop() {
			benchmarkPrivateResult = isPrivateStdlib(ip)
		}
	})
	b.Run("StdlibWithParseIP", func(b *testing.B) {
		for b.Loop() {
			benchmarkPrivateResult = isPrivateStdlib(net.ParseIP(benchmarkIPAddress))
		}
	})
}

func isPrivateStdlib(ip net.IP) bool {
	return ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast()
}

func legacyPrivateCIDRs() []*net.IPNet {
	blocks := []string{
		"127.0.0.1/8",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"169.254.0.0/16",
		"::1/128",
		"fc00::/7",
		"fe80::/10",
	}
	result := make([]*net.IPNet, 0, len(blocks))
	for _, block := range blocks {
		_, cidr, err := net.ParseCIDR(block)
		if err != nil {
			panic(err)
		}
		result = append(result, cidr)
	}
	return result
}

func isPrivateLegacy(address string, cidrs []*net.IPNet) bool {
	ip := net.ParseIP(address)
	for _, cidr := range cidrs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}
