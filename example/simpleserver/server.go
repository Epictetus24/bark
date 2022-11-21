package main

import (
	"net/http"

	"github.com/salukikit/bark/barkserv"
)

var (
	//these should always end in a trailing slash from the beacon standpoint followed by their ID
	reg  = []string{"/register/"}
	task = []string{"/task/"}
	out  = []string{"/out/"}
	ua   = "Mozilla/5.0 (Windows NT 10.0; Trident/7.0; rv:11.0) like Gecko"
)

// store your authconfig
var jwtconf *barkserv.JWTConf

func Register(w http.ResponseWriter, r *http.Request) {

	r.Write([]byte(ua))
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {

}

func OutputHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {

	jwtconf = barkserv.NewAuthConfig("supersecret", "HeaderName")

	myrouter := barkserv.NewBarkRouter(reg, task, out, jwtconf)

	myrouter.Regfunc = http.HandlerFunc(Register)

	barkhttpserver

}
