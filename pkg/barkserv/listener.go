package barkserv

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lucas-clemente/quic-go/http3"
)

type Server struct {
	Server     *http.Server
	Routes     *RouterConf
	TLS        *TLSConf
	Xi         *chi.Mux
	QuicServer *http3.Server
	Quic       bool
	LogPath    string
}

// Specifies the path for the barkweblog, please include full path with filename e.g. /path/to/logfile.log and not /path/to/
func SetWebLogger(path string) {

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		BarkLogger.Printf("Error creating web logfile: %v\n", err)
	}

	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(file, "[barkweb] ", log.Ldate|log.Ltime|log.LstdFlags)})
}

func (l *Server) StartListener(address string) {

	if l.LogPath != "" {
		SetWebLogger(l.LogPath)
	}

	l.Xi = chi.NewRouter()

	l.Xi.Use(middleware.RequestID)
	l.Xi.Use(middleware.RealIP)
	l.Xi.Use(middleware.Logger)
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
				BarkLogger.Printf("[barkserv] Error listening on %s\n", address)
				return
			}
		}()
		BarkLogger.Printf("[barkserv] Starting HTTP3(QUIC) Implant Server listening on %s\n", address)

		l.QuicServer = server

	} else if l.TLS != nil {

		go func() {
			err := httpserver.ListenAndServeTLS(l.TLS.Certpub, l.TLS.Certkey)
			if err != nil {
				BarkLogger.Printf("[barkserv] Error listening on %s\n", address)
				return
			}
		}()
		BarkLogger.Printf("[barkserv] Starting HTTPS Implant Server listening on %s\n", address)

	} else {
		go func() {
			err := httpserver.ListenAndServe()
			if err != nil {
				BarkLogger.Printf("[barkserv] Error listening on %s\n", address)
				return
			}
		}()
		BarkLogger.Printf("[barkserv] Starting HTTP Implant Server listening on %s\n", address)
	}

}

// Stop listeners
func (l *Server) StopListener() {

	if l.Quic == true {
		l.QuicServer.Close()

	} else {
		l.Server.Close()
	}

}
