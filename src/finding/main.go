package main

import (
	"fmt"
	"net"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type findingConf struct {
	Port string `default:"8001"`
}

func initXRay() {
	xray.Configure(xray.Config{})
}

func main() {
	var conf findingConf
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		appLogger.Fatal(err)
	}

	initXRay()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(xray.UnaryServerInterceptor()))
	findingServer := newFindingService() // DI service & repository
	finding.RegisterFindingServiceServer(server, findingServer)

	reflection.Register(server) // enable reflection API
	appLogger.Infof("Starting gRPC server at :%s", conf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Fatalf("Failed to gRPC serve: %v", err)
	}
}
