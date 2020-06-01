package main

import (
	"fmt"
	"net"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

type findingConf struct {
	Port string `default:"8081"`
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

	server := grpc.NewServer()
	findingServer := newFindingService(newFindingRepository()) // DI service & repository
	finding.RegisterFindingServiceServer(server, findingServer)
	appLogger.Infof("starting gRPC server at :%s", conf.Port)
	server.Serve(l)
}
