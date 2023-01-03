package barkserv

import (
	"log"
	"net/http"
	"os"
)

var (
	//Package wide logger for logging bark errors etc.
	BarkLogger *log.Logger
)

func init() {
	BarkLogger = log.New(os.Stdout, "[barkserv] ", log.Ldate|log.Ltime)
}

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
