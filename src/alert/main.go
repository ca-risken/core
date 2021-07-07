package main

import (
	"context"
	"fmt"
	"net"

	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/aws/aws-xray-sdk-go/xray"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type alertConf struct {
	Port    string `default:"8004"`
	EnvName string `default:"default" split_words:"true"`
}

func initXRay() {
	xray.Configure(xray.Config{
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
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				xray.UnaryServerInterceptor(),
				annotateEnvTracingInterceptor(conf.EnvName))))
	alertServer := newAlertService() // DI service & repository
	alert.RegisterAlertServiceServer(server, alertServer)

	reflection.Register(server) // enable reflection API
	appLogger.Infof("Starting gRPC server at :%s", conf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Fatalf("Failed to gRPC serve: %v", err)
	}
}

// TODO refactor
func annotateEnvTracingInterceptor(env string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := xray.AddAnnotation(ctx, "env", env); err != nil {
			appLogger.Warnf("failed to annotate environment to x-ray: %+v", err)
		}
		return handler(ctx, req)
	}
}
