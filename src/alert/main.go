package main

import (
	"fmt"
	"net"

	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type alertConf struct {
	Port string `default:"8004"`
}

func initXRay() {
	xray.Configure(xray.Config{
		DaemonAddr:     "127.0.0.1:2000",
		ServiceVersion: "TODO",
	})
}

func main() {
	var conf alertConf
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	// TODO toggle XRay
	initXRay()

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		appLogger.Fatal(err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(xray.UnaryServerInterceptor()))
	alertServer := newAlertService() // DI service & repository
	alert.RegisterAlertServiceServer(server, alertServer)

	reflection.Register(server) // enable reflection API
	appLogger.Infof("Starting gRPC server at :%s", conf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Fatalf("Failed to gRPC serve: %v", err)
	}
}
