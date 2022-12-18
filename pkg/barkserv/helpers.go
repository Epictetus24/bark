package barkserv

import (
	"crypto/tls"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)

var (
	dontverify = tls.Config{InsecureSkipVerify: true}
)

// Create a new BarkServ with HTTPS Defaults
func NewBarkServHTTPS(tlscert, tlskey string, rconf RouterConf) *HTTPSConf {
	httpsconf := &HTTPSConf{
		Server: &http.Server{},
		Routes: &rconf,
		TLS:    &TLSConf{Certpub: tlscert, Certkey: tlskey},
	}

	return httpsconf
}

// Create a new BarkServ with HTTP Defaults
func NewBarkServHTTP(rconf RouterConf) *HTTPSConf {
	httpsconf := &HTTPSConf{
		Server: &http.Server{},
		Routes: &rconf,
	}

	return httpsconf
}

// Create a new BarkServ router config with specified URIs. Handler funcs must be specified.
func NewBarkRouter(reguris, taskuris, outputuris []string) *RouterConf {
	rconf := &RouterConf{
		Reguris:  reguris,
		Taskuris: taskuris,
		Outuris:  outputuris,
	}

	return rconf
}

// Create a new BarkServ with QUIC defaults
func NewBarkServQUIC(tlscert, tlskey string, rconf RouterConf) *QUICConf {
	quicconf := &QUICConf{
		Server: &http3.Server{},
		Routes: &rconf,
		TLS:    &TLSConf{Certpub: tlscert, Certkey: tlskey},
	}

	//BarkServers[name] = quicconf
	return quicconf
}
