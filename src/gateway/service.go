package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

const (
	successJSONKey = "data"
	errorJSONKey   = "error"
)

type gatewayService struct {
	port          string
	findingClient finding.FindingServiceClient
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

	ctx := context.Background()
	conn, err := getGRPCConn(ctx, conf.FindingSvcAddr)
	if err != nil {
		return gatewayService{}, err
	}
	svc := gatewayService{
		port:          conf.Port,
		findingClient: finding.NewFindingServiceClient(conn),
	}
	return svc, nil
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		return nil, err
	}
	return conn, nil
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
