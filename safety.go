package bark

import (
	"crypto/tls"
	"fmt"
)

// GetTLSCertIssuer returns the issuer of the TLS certificate
// for the specified host and port.
func GetTLSCertIssuer(addr string, tcp bool) (string, error) {
	proto := "tcp"
	// Connect to the host and port
	if !tcp {
		proto = "udp"
	}
	conn, err := tls.Dial(proto, addr, nil)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Get the TLS certificate from the connection
	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return "", fmt.Errorf("no TLS certificate found")
	}

	// Return the issuer of the TLS certificate
	return certs[0].Issuer.String(), nil
}
