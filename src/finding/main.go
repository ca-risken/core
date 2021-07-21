package main

import (
	"fmt"
	"net"

	mimosaxray "github.com/CyberAgent/mimosa-common/pkg/xray"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/aws/aws-xray-sdk-go/xray"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type findingConf struct {
	Port    string `default:"8001"`
	EnvName string `default:"default" split_words:"true"`
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
