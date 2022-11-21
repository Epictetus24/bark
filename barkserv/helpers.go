package barkserv

import (
	"crypto/tls"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)

var (
	dontverify = tls.Config{InsecureSkipVerify: true}
)

// Create a new BarkServ with HTTP Defaults
func NewBarkServHTTPS(tlscert, tlskey string, rconf RouterConf) *HTTPSConf {
	httpsconf := &HTTPSConf{
		Server: &http.Server{},
		Routes: &rconf,
		TLS:    &TLSConf{Certpub: tlscert, Certkey: tlskey},
	}

	return httpsconf
}

// Create a new BarkServ router config with specified URIs. Handler funcs must be specified.
func NewBarkRouter(reguris, taskuris, outputuris []string, authconf *AuthConf) *RouterConf {
	rconf := &RouterConf{
		Reguris:  reguris,
		Taskuris: taskuris,
		Outuris:  outputuris,
	}

	return rconf
}

// Create a new BarkServ with QUIC defaults
func NewBarkServQUIC(tlscert, tlskey string, rconf RouterConf) *QUICConf {
	quicconf := &QUICConf{
		Server: &http3.Server{},
		Routes: &rconf,
		TLS:    &TLSConf{Certpub: tlscert, Certkey: tlskey},
	}

	//BarkServers[name] = quicconf
	return quicconf
}

func NewAuthConfig(secret, name string) *AuthConf {
	authconf := &AuthConf{}
	authconf.Secret = secret
	authconf.Name = name
	authconf.TokenAuth = newJWTAuth(authconf.Secret)

	return authconf

}

func GetBarkJWTCookie(r *http.Request, cookiename string) string {
	cookie, err := r.Cookie(cookiename)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (ac *AuthConf) NewBarkCookie(domain, cookiename, id, value string, secure, httponly bool) http.Cookie {
	barktoken := newBarkJWT(id, value, *ac.TokenAuth)
	barkcookie := http.Cookie{
		Name:  cookiename,
		Value: barktoken,

		Domain:   domain,
		HttpOnly: httponly,
		Secure:   secure,
	}

	return barkcookie

}

func (ac *AuthConf) NewBarkHeader(cookiename, id, value string) http.Cookie {
	barktoken := newBarkJWT(id, value, *ac.TokenAuth)
	barkcookie := http.Cookie{
		Name:  cookiename,
		Value: barktoken,

		Domain:   domain,
		HttpOnly: httponly,
		Secure:   secure,
	}

	return barkcookie

}
