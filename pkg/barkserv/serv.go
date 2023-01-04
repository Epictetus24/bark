package barkserv

import (
	"github.com/go-chi/chi/v5"
)

func NewRoutes(xi *chi.Mux, rconf *RouterConf) {

	xi.Group(func(xi chi.Router) {
		// Seek, verify and validate JWT tokens
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
