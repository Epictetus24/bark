package main

import (
	"encoding/base64"
	"net/http"

	"github.com/salukikit/bark/barkserv"
)

var (
	//these should always end in a trailing slash from the beacon standpoint followed by their ID
	reg  = []string{"/register"}
	task = []string{"/task"}
	out  = []string{"/out"}
	ua   = "Mozilla/5.0 (Windows NT 10.0; Trident/7.0; rv:11.0) like Gecko"
)

func init() {

}

// Register Bark Implants
func Register(w http.ResponseWriter, r *http.Request) {

	// "Encrypt" the body (new user agent) and send
	var body []byte
	base64.StdEncoding.Encode(body, []byte(ua))
	w.Write(body)

}

// Handle tasks
func TaskHandler(w http.ResponseWriter, r *http.Request) {

}

// Handle output sent by implant
func OutputHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {

	//create our router
	myrouter := barkserv.NewBarkRouter(reg, task, out)

	myrouter.Regfunc = http.HandlerFunc(Register)

	barkhttpserver := barkserv.NewBarkServHTTPS("example.crt", "example.key", *myrouter)

	barkhttpserver.StartListener("127.0.0.1:8080", true)

	x := true

	for x {

		/* wait for some reason until you wanna stop the server do something else
		Can be stopped with barkhttpserver.StopListener()

		*/

	}

}
