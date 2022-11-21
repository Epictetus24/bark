package barkserv

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lucas-clemente/quic-go/http3"
)

type QUICConf struct {
	Routes *RouterConf
	TLS    *TLSConf
	Server *http3.Server
	Xi     *chi.Mux
}

func (qc *QUICConf) StartListener(address string, verbose bool) *http3.Server {

	qc.Xi = chi.NewRouter()
	if verbose {
		qc.Xi.Use(middleware.Logger)
	}
	qc.Xi.Use(middleware.RequestID)
	qc.Xi.Use(middleware.RealIP)
	qc.Xi.Use(middleware.Recoverer)
	qc.Xi.Use(middleware.Recoverer)
	qc.Xi.Use(jwtauth.Verifier(qc.Routes.AuthConf.TokenAuth))

	newRoutes(qc.Xi, qc.Routes, qc.Routes.AuthConf)

	server := &http3.Server{
		Addr:    address,
		Handler: qc.Xi,
	}

	go func() {
		err := server.ListenAndServeTLS(qc.TLS.Certpub, qc.TLS.Certkey)
		if err != nil {
			log.Fatal(err)
		}
	}()

	qc.Server = server

	return server

}

func (qc *QUICConf) StopListener() {

	qc.Server.Close()

}
