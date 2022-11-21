package main

import "github.com/salukikit/bark/barkserv"

func main() {


	/*
	type ListenerConf struct {
	//URLS to register
	Taskurls []string
	Taskfunc http.HandlerFunc
	Outurls  []string
	Outfunc  http.HandlerFunc
	Regurls  []string
	Regfunc  http.HandlerFunc

	//general settings
	Certpub  string
	Certkey  string
	Verbose  bool
	Bindaddr string
	Bindport string

	//Auth settings
	Jwt bool
	Pki bool
	E2e bool
}
	*/

	lconf := &barkserv.ListenerConf{
		Certpub: ,


		Jwt: true,

		Taskurls: []string{"/task"},

	}

	barkserv.StartHTTPSBarkListener(lconf)

}
