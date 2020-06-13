package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func newRouter(svc gatewayService) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(httpLogger)
	r.Use(middleware.StripSlashes) // convert URI path. like `/hoge/111/` -> `/hoge/111`

	r.Route("/finding", func(r chi.Router) {
		r.Get("/", svc.listFindingHandler)
		r.Get("/{finding_id}", svc.getFindingHandler)
		r.Get("/{finding_id}/tag", svc.listFindingTagHandler)
		r.Post("/put", svc.putFindingHandler)
		r.Post("/delete", svc.deleteFindingHandler)
		r.Post("/tag", svc.tagFindingHandler)
		r.Post("/untag", svc.untagFindingHandler)
	})
	r.Route("/resource", func(r chi.Router) {
		r.Get("/", svc.listResourceHandler)
		r.Get("/{resource_id}", svc.getResourceHandler)
		r.Post("/put", svc.putResourceHandler)
	})

	return r
}
