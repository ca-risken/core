package main

import (
	"net/http"

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
		r.Get("/detail", svc.getFindingHandler)
		r.Get("/tag", svc.listFindingTagHandler)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))
			r.Post("/put", svc.putFindingHandler)
			r.Post("/delete", svc.deleteFindingHandler)
			r.Post("/tag/put", svc.tagFindingHandler)
			r.Post("/tag/delete", svc.untagFindingHandler)
		})
	})
	r.Route("/resource", func(r chi.Router) {
		r.Get("/", svc.listResourceHandler)
		r.Get("/detail", svc.getResourceHandler)
		r.Get("/tag", svc.listResourceTagHandler)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AllowContentType("application/json"))
			r.Post("/put", svc.putResourceHandler)
			r.Post("/delete", svc.deleteResourceHandler)
			r.Post("/tag/put", svc.tagResourceHandler)
			r.Post("/tag/delete", svc.untagResourceHandler)
		})
	})

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	return r
}
