package netutil

import (
	"net"

	"github.com/cockroachdb/errors"
	"github.com/effective-security/xlog"
)

var logger = xlog.NewPackageLogger("github.com/effective-security/x", "netutil")

// FindFreePort returns a free port found on a host
func FindFreePort(host string, maxAttempts int) (int, error) {
	if host == "" {
		host = "localhost"
	}
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			lastErr = errors.Wrapf(err, "unable to resolve TCP address for host %s", host)
			logger.KV(xlog.ERROR,
				"reason", "unable to resolve tcp addr",
				"err", err.Error())
			continue
		}
		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			lastErr = errors.Wrapf(err, "unable to listen on %s", addr)
			logger.KV(xlog.ERROR,
				"reason", "unable to listen",
				"addr", addr,
				"err", err.Error())
			continue
		}

		port := l.Addr().(*net.TCPAddr).Port
		if err := l.Close(); err != nil {
			return 0, errors.Wrapf(err, "unable to release TCP port %d", port)
		}
		return port, nil
	}

	return 0, errors.WithMessage(lastErr, "no free port found")
}
