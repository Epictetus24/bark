package barkserv

import (
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
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
	AuthConf *AuthConf
}

type AuthConf struct {
	//Auth settings
	TokenAuth *jwtauth.JWTAuth //JWT Context
	Name      string           //token name
	Secret    string
	RedirUrl  string //url for redirecting unauthed targets
	Expiry    time.Duration
}
