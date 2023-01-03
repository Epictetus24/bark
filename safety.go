package bark

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
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

type AuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expires     string `json:"expires_at"`
}

func BuryinJwt(data []byte) (string, error) {
	var auth AuthToken
	fakeexpiry := time.Now().Add(time.Minute * 10)
	auth.AccessToken = base64.StdEncoding.EncodeToString(data)
	auth.TokenType = "bearer"
	auth.Expires = fakeexpiry.String()

	b, err := json.MarshalIndent(auth, "", "  ")
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(b), nil
}
