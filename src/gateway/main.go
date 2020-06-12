package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kelseyhightower/envconfig"
)

type gatewayConf struct {
	Port           string `default:"8080"`
	FindingSvcAddr string `required:"true" split_words:"true"`
}

func main() {
	var conf gatewayConf
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	ctx := context.Background()
	svc, err := newGatewayService(ctx, conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	r := chi.NewRouter()
	r.Use(middleware.StripSlashes) // convert URI path. like `/hoge/111/` -> `/hoge/111`
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(httpLogger)
	r.Get("/finding", svc.listFindingHandler)
	r.Get("/finding/{finding_id}", svc.getFindingHandler)
	r.Post("/finding/put", svc.putFindingHandler)
	r.Post("/finding/delete", svc.deleteFindingHandler)
	r.Get("/finding/{finding_id}/tag", svc.listFindingTagHandler)
	r.Post("/finding/tag", svc.tagFindingHandler)

	appLogger.Infof("starting http server at :%s", conf.Port)
	err = http.ListenAndServe(":"+conf.Port, r)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
}
