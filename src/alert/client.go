package main

import (
	"context"
	"time"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

type findingConfig struct {
	FindingSvcAddr string `required:"true" split_words:"true"`
}

func newFindingClient() finding.FindingServiceClient {
	var conf findingConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatalf("Faild to load finding config error: err=%+v", err)
	}

	ctx := context.Background()
	conn, err := getGRPCConn(ctx, conf.FindingSvcAddr)
	if err != nil {
		appLogger.Fatalf("Faild to get GRPC connection: err=%+v", err)
	}
	return finding.NewFindingServiceClient(conn)
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
