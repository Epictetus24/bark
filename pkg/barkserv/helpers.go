package barkserv

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/salukikit/bark"
)

var (
	dontverify = tls.Config{InsecureSkipVerify: true}
)

// Create a new BarkServ with HTTPS Defaults
func NewBarkServHTTPS(tlscert, tlskey string, rconf RouterConf) Server {

	return Server{
		Server: &http.Server{},
		Routes: &rconf,
		TLS:    &TLSConf{Certpub: tlscert, Certkey: tlskey},
	}
}

// Create a new BarkServ with HTTP Defaults
func NewBarkServHTTP(rconf RouterConf) Server {

	return Server{
		Server: &http.Server{},
		Routes: &rconf,
	}
}

// Create a new BarkServ router config with specified URIs. Handler funcs must be specified.
func NewBarkRouter(reguris, taskuris, outputuris []string) RouterConf {

	return RouterConf{
		Reguris:  reguris,
		Taskuris: taskuris,
		Outuris:  outputuris,
	}

}

// Create a new BarkServ with QUIC defaults
func NewBarkServQUIC(tlscert, tlskey string, rconf RouterConf) Server {

	return Server{
		QuicServer: &http3.Server{},
		Routes:     &rconf,
		TLS:        &TLSConf{Certpub: tlscert, Certkey: tlskey},
		Quic:       true,
	}
}

// Strips the fakestuff off the Auth header JWT.
// Pass the authheader data from a barkmessage and get back the encrypted access token as []byte.
func DataFromFakeJwt(authHead string) ([]byte, error) {
	var fakeauth *bark.AuthToken
	if authHead == "" {
		return nil, fmt.Errorf("Empty string, no auth header")
	}

	fakestr := strings.ReplaceAll(authHead, "Bearer ", "")
	fakedat, err := base64.StdEncoding.DecodeString(fakestr)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(fakedat, fakeauth)
	if err != nil {
		return nil, err
	}
	encdat, err := base64.StdEncoding.DecodeString(fakeauth.AccessToken)
	if err != nil {
		return nil, err
	}

	return encdat, nil

}
