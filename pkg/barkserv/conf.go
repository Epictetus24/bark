package barkserv

import (
	"net/http"
)

type TLSConf struct {
	//TLS settings
	Certpub string
	Certkey string
}

type RouterConf struct {
	//URLS to register & their respective handlers
	Taskuris []string
	Taskfunc http.HandlerFunc
	Outuris  []string
	Outfunc  http.HandlerFunc
	Reguris  []string
	Regfunc  http.HandlerFunc
}
