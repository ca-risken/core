package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
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
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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
	MaxAnalyzeAPICall int64
	BaseURL           string
}

func NewConfig(maxAnalyzeAPICall int64, baseURL string) Config {
	return Config{
		MaxAnalyzeAPICall: maxAnalyzeAPICall,
		BaseURL:           baseURL,
	}
}

func (s *Server) Run(ctx context.Context) error {
	clientAddr := fmt.Sprintf("localhost:%s", s.port)
	fc, err := s.newFindingClient(clientAddr)
	if err != nil {
		return err
	}
	pc, err := s.newProjectClient(clientAddr)
	if err != nil {
		return err
	}
	iamc, err := s.newIAMClient(clientAddr)
	if err != nil {
		return err
	}
	isvc := iamserver.NewIAMService(s.db, fc, s.logger)
	asvc := alertserver.NewAlertService(
		s.config.MaxAnalyzeAPICall,
		s.config.BaseURL,
		fc,
		pc,
		s.db,
		s.logger,
	)
	fsvc := findingserver.NewFindingService(s.db, s.logger)
	psvc := projectserver.NewProjectService(s.db, iamc, s.logger)
	rsvc := reportserver.NewReportService(s.db, s.logger)
	hsvc := health.NewServer()

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				grpctrace.UnaryServerInterceptor(),
				mimosarpc.LoggingUnaryServerInterceptor(s.logger))))
	iam.RegisterIAMServiceServer(server, isvc)
	report.RegisterReportServiceServer(server, rsvc)
	alert.RegisterAlertServiceServer(server, asvc)
	finding.RegisterFindingServiceServer(server, fsvc)
	project.RegisterProjectServiceServer(server, psvc)
	grpc_health_v1.RegisterHealthServer(server, hsvc)

	reflection.Register(server) // enable reflection API

	s.logger.Infof(ctx, "Starting gRPC server at :%s", s.port)
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	errChan := make(chan error)
	go func() {
		if err := server.Serve(l); err != nil && err != grpc.ErrServerStopped {
			s.logger.Errorf(ctx, "failed to serve grpc: %w", err)
			errChan <- err
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := healthCheck(ctx, clientAddr); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			s.logger.Errorf(ctx, "health check is failed: %w", err)
		} else {
			fmt.Fprintln(w, "ok")
		}
	})

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%s:3000", s.host), mux); err != http.ErrServerClosed {
			s.logger.Errorf(ctx, "failed to start http server: %w", err)
			errChan <- err
		}
	}()

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		s.logger.Info(ctx, "Shutdown gRPC server...")
		server.GracefulStop()
	}

	return nil
}

func healthCheck(ctx context.Context, addr string) error {
	conn, err := getGRPCConn(context.Background(), addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	res, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return err
	}
	if res.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return fmt.Errorf("returned status is '%v'", res.Status)
	}

	return nil
}

func (s *Server) newFindingClient(svcAddr string) (finding.FindingServiceClient, error) {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc connection: err=%w", err)
	}
	return finding.NewFindingServiceClient(conn), nil
}

func (s *Server) newIAMClient(svcAddr string) (iam.IAMServiceClient, error) {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc connection: err=%w", err)
	}
	return iam.NewIAMServiceClient(conn), nil
}

func (s *Server) newProjectClient(svcAddr string) (project.ProjectServiceClient, error) {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc connection: err=%w", err)
	}
	return project.NewProjectServiceClient(conn), nil
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
