package bark

import (
	"crypto/tls"
	"net/http"
	"reflect"

	"github.com/lucas-clemente/quic-go/http3"
)

var (
	dontverify = tls.Config{InsecureSkipVerify: true}
)

// Create a new Barker with HTTP Defaults
func NewBarkerHTTP(name string, verifytls bool) *BarkConfig {
	httpconf := &BarkConfig{}
	if verifytls {
		httpconf.tr = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		}
	} else {
		httpconf.tr = &http.Transport{
			TLSClientConfig: &dontverify,
			Proxy:           http.ProxyFromEnvironment,
		}

	}
	Barkers[name] = httpconf
	return httpconf
}

// Create a new Barker with QUIC defaults
func NewBarkerQUIC(name string, verifytls bool) *BarkConfig {
	httpconf := &BarkConfig{}
	if verifytls {
		httpconf.tr = &http3.RoundTripper{}
	} else {
		httpconf.tr = &http3.RoundTripper{
			TLSClientConfig: &dontverify,
		}

	}
	Barkers[name] = httpconf
	return httpconf
}

// Create a new barker with no settings
func NewBarker(name string) *BarkConfig {
	httpconf := &BarkConfig{}
	Barkers[name] = httpconf
	return httpconf
}

// UpdateBarker updates the barker's comms config with the provided BarkConfig
func UpdateBarkers(barkername string, httpconf *BarkConfig) {
	Barkers[barkername] = httpconf
}

// Delete an existing Barker
func DeleteBarker(name string) {
	delete(Barkers, name)
}

func ListBarkersName() []string {
	keys := []string{}
	value := reflect.ValueOf(Barkers)
	if value.Kind() == reflect.Map {
		for _, v := range value.MapKeys() {
			if v.Kind() == reflect.String {
				keys = append(keys, v.String())
			}
		}
	}
	return keys

}

func ListAllBarkers() []Barker {
	var allbarkers []Barker
	for _, v := range Barkers {
		allbarkers = append(allbarkers, v)

	}
	return allbarkers

}
