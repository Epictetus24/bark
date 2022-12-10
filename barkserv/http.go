package barkserv

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTPSConf struct {
	Server *http.Server
	Routes *RouterConf
	TLS    *TLSConf
	Xi     *chi.Mux
}

type HTTPConf struct {
	Server *http.Server
	Routes *RouterConf
	Xi     *chi.Mux
}

func (hc *HTTPConf) StartListener(address string, verbose bool) *http.Server {

	hc.Xi = chi.NewRouter()
	hc.Xi.Use(middleware.RequestID)
	hc.Xi.Use(middleware.RealIP)
	if verbose == true {
		hc.Xi.Use(middleware.Logger)
	}
	hc.Xi.Use(middleware.Recoverer)

	newRoutes(hc.Xi, hc.Routes)

	color.Green("Starting HTTP Implant Server listening on %s\n", address)

	server := &http.Server{
		Addr:    address,
		Handler: hc.Xi,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	hc.Server = server
	return server

}

func (hcs *HTTPSConf) StartListener(address string, verbose bool) *http.Server {
	hcs.Xi = chi.NewRouter()
	hcs.Xi.Use(middleware.RequestID)
	hcs.Xi.Use(middleware.RealIP)
	if verbose == true {
		hcs.Xi.Use(middleware.Logger)
	}
	hcs.Xi.Use(middleware.Recoverer)

	newRoutes(hcs.Xi, hcs.Routes)

	color.Green("Starting HTTPS Implant Server listening on %s\n", address)

	server := &http.Server{
		Addr:    address,
		Handler: hcs.Xi,
	}

	go func() {
		err := server.ListenAndServeTLS(hcs.TLS.Certpub, hcs.TLS.Certkey)
		if err != nil {
			log.Fatal(err)
		}
	}()

	hcs.Server = server

	return server

}

func (hcs *HTTPSConf) StopListener() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hcs.Server.Shutdown(ctx)

}

func (hc *HTTPConf) StopListener() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	hc.Server.Shutdown(ctx)

}
