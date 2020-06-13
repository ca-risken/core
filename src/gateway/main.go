package main

import (
	"net/http"
)

func main() {
	svc, err := newGatewayService()
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	appLogger.Infof("starting http server at :%s", svc.port)
	err = http.ListenAndServe(":"+svc.port, newRouter(svc))
	if err != nil {
		appLogger.Fatal(err.Error())
	}
}
