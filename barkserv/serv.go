package barkserv

import (
	"net/http"

	"github.com/fatih/color"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/lucas-clemente/quic-go/http3"
)

func NewRoutes(xi *chi.Mux, lconf ListenerConf) {

	//register task routes
	for i := range lconf.Taskurls {
		taskpath := lconf.Taskurls[i] + "{impid}"
		xi.Get(taskpath, lconf.Taskfunc)

	}

	//register output routes
	for i := range lconf.Outurls {
		outpath := lconf.Outurls[i] + "{impid}"
		xi.Post(outpath, lconf.Outfunc)

	}

	for i := range lconf.Regurls {
		regpath := lconf.Regurls[i] + "{impid}"
		xi.Get(regpath, lconf.Regfunc)

	}

}

// Start HTTPS REST Server with Listenerconfig
func StartHTTPSBarkListener(lconf ListenerConf) {
	restHandleRequests(lconf, 1)

}

// Start QUIC REST Server with Listenerconfig
func StartQUICBarkListener(lconf ListenerConf) {
	restHandleRequests(lconf, 2)

}

// Not Recommended: Start HTTP REST Server with Listenerconfig
func StartHTTPBarkListener(lconf ListenerConf) {
	restHandleRequests(lconf, 3)

}

// Rest Handle requests handles routing of API requests
func restHandleRequests(lconf ListenerConf, proto int) {
	xi := chi.NewRouter()
	if lconf.Verbose {
		xi.Use(middleware.Logger)
	}
	xi.Use(middleware.RequestID)
	xi.Use(middleware.RealIP)
	xi.Use(middleware.Recoverer)

	NewRoutes(xi, lconf)

	address := lconf.Bindaddr + ":" + lconf.Bindport

	switch proto {
	case 2:
		go startQUICListener(address, lconf.Certpub, lconf.Certkey, xi)

	case 3:
		go startHTTPListener(address, xi)

	default:
		go startHTTPSListener(address, lconf.Certpub, lconf.Certkey, xi)

	}

}

func startHTTPListener(address string, xi *chi.Mux) error {

	color.Green("Starting HTTP(s) Implant Server listening on %s\n", address)

	err := http.ListenAndServe(address, xi)
	return err
}

func startHTTPSListener(address, cert, key string, xi *chi.Mux) error {

	color.Green("Starting HTTP(s) Implant Server listening on %s\n", address)

	err := http.ListenAndServeTLS(address, cert, key, xi)

	return err

}

func startQUICListener(address, cert, key string, xi *chi.Mux) error {

	color.Green("Starting QUIC Implant Server listening on %s\n", address)

	err := http3.ListenAndServeQUIC(address, cert, key, xi)
	return err

}
