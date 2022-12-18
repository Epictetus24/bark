package barkserv

import (
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucas-clemente/quic-go/http3"
)

type Listen struct {
	Server     *http.Server
	Routes     *RouterConf
	TLS        *TLSConf
	Xi         *chi.Mux
	QuicServer *http3.Server
	Quic       bool
}

func (l *Listen) StartListener(address string, verbose bool) {

	l.Xi = chi.NewRouter()
	if verbose {
		l.Xi.Use(middleware.Logger)
	}
	l.Xi.Use(middleware.RequestID)
	l.Xi.Use(middleware.RealIP)
	l.Xi.Use(middleware.Recoverer)

	newRoutes(l.Xi, l.Routes)

	httpserver := &http.Server{
		Addr:    address,
		Handler: l.Xi,
	}

	if l.Quic {

		server := &http3.Server{
			Addr:    address,
			Handler: l.Xi,
		}
		go func() {
			err := server.ListenAndServeTLS(l.TLS.Certpub, l.TLS.Certkey)
			if err != nil {
				color.Red("[barkserv] Error listening on %s\n", address)
				log.Fatal(err)
			}
		}()
		color.Green("[barkserv] Starting HTTP3(QUIC) Implant Server listening on %s\n", address)

		l.QuicServer = server

	} else if l.TLS != nil {

		go func() {
			err := httpserver.ListenAndServeTLS(l.TLS.Certpub, l.TLS.Certkey)
			if err != nil {
				color.Red("[barkserv] Error listening on %s\n", address)
				log.Fatal(err)
			}
		}()
		color.Green("[barkserv] Starting HTTPS Implant Server listening on %s\n", address)

	} else {
		go func() {
			err := httpserver.ListenAndServe()
			if err != nil {
				color.Red("[barkserv] Error listening on %s\n", address)
				log.Fatal(err)
			}
		}()
		color.Green("[barkserv] Starting HTTP Implant Server listening on %s\n", address)
	}

}

// Stop listeners
func (l *Listen) StopListener() {

	if l.Quic == true {
		l.QuicServer.Close()

	} else {
		l.Server.Close()
	}

}
