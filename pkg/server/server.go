package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	mimosarpc "github.com/ca-risken/common/pkg/rpc"
	"github.com/ca-risken/core/pkg/db"
	alertserver "github.com/ca-risken/core/pkg/server/alert"
	findingserver "github.com/ca-risken/core/pkg/server/finding"
	iamserver "github.com/ca-risken/core/pkg/server/iam"
	projectserver "github.com/ca-risken/core/pkg/server/project"
	reportserver "github.com/ca-risken/core/pkg/server/report"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/core/proto/report"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type Server struct {
	host   string
	port   string
	db     *db.Client
	logger logging.Logger
	config Config
}

func NewServer(host string, port string, db *db.Client, logger logging.Logger, config Config) *Server {
	return &Server{
		host:   host,
		port:   port,
		db:     db,
		logger: logger,
		config: config,
	}
}

type Config struct {
	MaxAnalyzeAPICall    int64
	NotificationAlertURL string
}

func NewConfig(maxAnalyzeAPICall int64, notificationAlertURL string) Config {
	return Config{
		MaxAnalyzeAPICall:    maxAnalyzeAPICall,
		NotificationAlertURL: notificationAlertURL,
	}
}

func (s *Server) Run() error {
	clientAddr := fmt.Sprintf("localhost:%s", s.port)
	fc := s.newFindingClient(clientAddr)
	isvc := iamserver.NewIAMService(s.db, fc)
	asvc := alertserver.NewAlertService(
		s.config.MaxAnalyzeAPICall,
		s.config.NotificationAlertURL,
		fc,
		s.newProjectClient(clientAddr),
		s.db,
	)
	fsvc := findingserver.NewFindingService(s.db)
	psvc := projectserver.NewProjectService(s.db, s.newIAMClient(clientAddr))
	rsvc := reportserver.NewReportService(s.db)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				mimosarpc.LoggingUnaryServerInterceptor(s.logger),
				grpctrace.UnaryServerInterceptor())))
	iam.RegisterIAMServiceServer(server, isvc)
	report.RegisterReportServiceServer(server, rsvc)
	alert.RegisterAlertServiceServer(server, asvc)
	finding.RegisterFindingServiceServer(server, fsvc)
	project.RegisterProjectServiceServer(server, psvc)

	reflection.Register(server) // enable reflection API

	s.logger.Infof("Starting gRPC server at :%s", s.port)
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	errChan := make(chan error)
	go func() {
		if err := server.Serve(l); err != nil && err != grpc.ErrServerStopped {
			errChan <- fmt.Errorf("failed to serve: %w", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		s.logger.Info("Shutdown gRPC server...")
		server.GracefulStop()
	}

	return nil
}

func (s *Server) newFindingClient(svcAddr string) finding.FindingServiceClient {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		s.logger.Fatalf("failed to get grpc connection: err=%+v", err)
	}
	return finding.NewFindingServiceClient(conn)
}

func (s *Server) newIAMClient(svcAddr string) iam.IAMServiceClient {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		s.logger.Fatalf("failed to get grpc connection: err=%+v", err)
	}
	return iam.NewIAMServiceClient(conn)
}

func (s *Server) newProjectClient(svcAddr string) project.ProjectServiceClient {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		s.logger.Fatalf("failed to get grpc connection: err=%+v", err)
	}
	return project.NewProjectServiceClient(conn)
}

func getGRPCConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithUnaryInterceptor(grpctrace.UnaryClientInterceptor()),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
