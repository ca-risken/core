package main

import (
	"fmt"
	"net"

	"github.com/aws/aws-xray-sdk-go/xray"
	mimosarpc "github.com/ca-risken/common/pkg/rpc"
	mimosaxray "github.com/ca-risken/common/pkg/xray"
	"github.com/ca-risken/core/proto/finding"
	"github.com/gassara-kys/envconfig"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type findingConf struct {
	Port    string `default:"8001"`
	EnvName string `default:"local" split_words:"true"`
}

func main() {
	var conf findingConf
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	mimosaxray.InitXRay(xray.Config{})

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		appLogger.Fatal(err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				mimosarpc.LoggingUnaryServerInterceptor(appLogger),
				xray.UnaryServerInterceptor(),
				mimosaxray.AnnotateEnvTracingUnaryServerInterceptor(conf.EnvName))))
	findingServer := newFindingService() // DI service & repository
	finding.RegisterFindingServiceServer(server, findingServer)

	reflection.Register(server) // enable reflection API
	appLogger.Infof("Starting gRPC server at :%s", conf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Fatalf("Failed to gRPC serve: %v", err)
	}
}
