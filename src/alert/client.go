package main

import (
	"context"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/project"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newFindingClient(svcAddr string) finding.FindingServiceClient {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		appLogger.Fatalf("Faild to get GRPC connection: err=%+v", err)
	}
	return finding.NewFindingServiceClient(conn)
}

func newProjectClient(svcAddr string) project.ProjectServiceClient {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		appLogger.Fatalf("Faild to get GRPC connection: err=%+v", err)
	}
	return project.NewProjectServiceClient(conn)
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithUnaryInterceptor(xray.UnaryClientInterceptor()),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
