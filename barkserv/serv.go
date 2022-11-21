package barkserv

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

var nulltime time.Duration

func newRoutes(xi *chi.Mux, rconf *RouterConf, jwt *AuthConf) {

	xi.Group(func(xi chi.Router) {
		// Seek, verify and validate JWT tokens
		xi.Use(jwtauth.Verifier(jwt.TokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware.
		xi.Use(jwtauth.Authenticator)

		for i := range rconf.Taskuris {
			taskpath := rconf.Taskuris[i]
			xi.Get(taskpath, rconf.Taskfunc)

		}
		//register output routes
		for i := range rconf.Outuris {
			outpath := rconf.Outuris[i]
			xi.Post(outpath, rconf.Outfunc)

		}

	})

	//register registration routes
	for i := range rconf.Reguris {
		regpath := rconf.Reguris[i]
		xi.Get(regpath, rconf.Regfunc)

	}

}

func newJWTAuth(secret string) *jwtauth.JWTAuth {
	tokenAuth := jwtauth.New("HS256", []byte(secret), nil)

	return tokenAuth
}

func newBarkJWT(id, value string, ac *AuthConf) string {
	barktoken := make(map[string]interface{})
	barktoken["id"] = id
	barktoken["value"] = value
	jwtauth.SetIssuedNow(barktoken)

	if ac.Expiry > nulltime {
		jwtauth.ExpireIn(ac.Expiry)
	}

	_, tokenString, _ := ac.TokenAuth.Encode(barktoken)
	return tokenString
}
