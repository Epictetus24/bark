package barkserv

import (
	"log"
	"net/http"
)

var (
	//Package wide logger for logging bark errors etc.
	BarkLogger *log.Logger
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
