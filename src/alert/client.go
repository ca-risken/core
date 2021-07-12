package main

import (
	"context"
	"time"

	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-core/proto/project"
	"github.com/aws/aws-xray-sdk-go/xray"
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

type projectConfig struct {
	ProjectSvcAddr string `required:"true" split_words:"true"`
}

func newProjectClient() project.ProjectServiceClient {
	var conf projectConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatalf("Faild to load project config error: err=%+v", err)
	}

	ctx := context.Background()
	conn, err := getGRPCConn(ctx, conf.ProjectSvcAddr)
	if err != nil {
		appLogger.Fatalf("Faild to get GRPC connection: err=%+v", err)
	}
	return project.NewProjectServiceClient(conn)
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	// TODO refactoring xray
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithUnaryInterceptor(xray.UnaryClientInterceptor()), grpc.WithInsecure(), grpc.WithTimeout(time.Second*3))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
