package bark

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"time"

	"github.com/lucas-clemente/quic-go/http3"
)

var (
	dontverify = tls.Config{InsecureSkipVerify: true}
)

func Jitter(d time.Duration, j float64) time.Duration {
	if j < 0.0 {
		return d
	}

	r := rand.Float64() * float64(d)
	if j > 0.0 && j < 1.0 {
		r = float64(j)*r + float64(1.0-j)*float64(d)
	}

	return time.Duration(r)
}

// Create a new Barker with HTTP Defaults
func NewBarkerHTTP(name string, verifytls bool) *BarkConfig {
	httpconf := &BarkConfig{}
	if verifytls {
		httpconf.Tr = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		}
	} else {
		httpconf.Tr = &http.Transport{
			TLSClientConfig: &dontverify,
			Proxy:           http.ProxyFromEnvironment,
		}

	}

	return httpconf
}

// Create a new Barker with QUIC defaults
func NewBarkerQUIC(name string, verifytls bool) *BarkConfig {
	httpconf := &BarkConfig{}
	if verifytls {
		httpconf.Tr = &http3.RoundTripper{}
	} else {
		httpconf.Tr = &http3.RoundTripper{
			TLSClientConfig: &dontverify,
		}

	}

	return httpconf
}

// Create a new barker with no settings
func NewBarker(name string) *BarkConfig {
	httpconf := &BarkConfig{}

	return httpconf
}
