package main

import (
	"context"
	"net/http"

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

	r := newRouoter(svc)
	appLogger.Infof("starting http server at :%s", conf.Port)
	err = http.ListenAndServe(":"+conf.Port, r)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
}
