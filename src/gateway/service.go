package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/gddo/httputil/header"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

type gatewayService struct {
	port           string
	findingSvcAddr string
	findingSvcConn *grpc.ClientConn
}

type gatewayConf struct {
	Port           string `default:"8080"`
	FindingSvcAddr string `required:"true" split_words:"true"`
}

func newGatewayService() (gatewayService, error) {
	var conf gatewayConf
	err := envconfig.Process("", &conf)
	if err != nil {
		return gatewayService{}, err
	}

	svc := gatewayService{
		port:           conf.Port,
		findingSvcAddr: conf.FindingSvcAddr,
	}
	ctx := context.Background()
	if err := mustConnGRPC(ctx, &svc.findingSvcConn, svc.findingSvcAddr); err != nil {
		return svc, err
	}
	return svc, nil
}

func mustConnGRPC(ctx context.Context, conn **grpc.ClientConn, addr string) error {
	var err error
	*conn, err = grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Second*3),
	)
	if err != nil {
		return err
	}
	return nil
}

func writeResponse(w http.ResponseWriter, status int, body map[string]interface{}) {
	if body == nil {
		w.WriteHeader(status)
		return
	}
	buf, err := json.Marshal(body)
	if err != nil {
		appLogger.Errorf("Response body JSON marshal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	w.Write(buf)
}

func validatePostHeader(r *http.Request) error {
	if r.Header.Get("Content-Type") == "" {
		return errors.New("Not found Content-Type header")
	}
	value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
	if value != "application/json" {
		return fmt.Errorf("Unexpected Content-Type. want=application/json, got=%s", value)
	}
	return nil
}
