package serv

import "net/http"

type ListenerConf struct {
	//URLS to register
	Taskurls []string
	Taskfunc http.HandlerFunc
	Outurls  []string
	Outfunc  http.HandlerFunc
	Regurls  []string
	Regfunc  http.HandlerFunc

	//general settings
	Proto    int
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
