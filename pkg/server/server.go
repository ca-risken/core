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
	aiserver "github.com/ca-risken/core/pkg/server/ai"
	alertserver "github.com/ca-risken/core/pkg/server/alert"
	findingserver "github.com/ca-risken/core/pkg/server/finding"
	iamserver "github.com/ca-risken/core/pkg/server/iam"
	organizationserver "github.com/ca-risken/core/pkg/server/organization"
	organization_iamserver "github.com/ca-risken/core/pkg/server/organization_iam"
	projectserver "github.com/ca-risken/core/pkg/server/project"
	reportserver "github.com/ca-risken/core/pkg/server/report"
	"github.com/ca-risken/core/proto/ai"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/organization_iam"
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
	MaxAnalyzeAPICall       int64
	BaseURL                 string
	OpenAIToken             string
	ChatGPTModel            string
	ReasoningModel          string
	defaultLocale           string
	excludeDeleteDataSource []string
	SlackAPIToken           string
}

func NewConfig(maxAnalyzeAPICall int64, baseURL, openaiToken, chatGPTModel, reasoningModel, defaultLocale, SlackAPIToken string, excludeDeleteDataSource []string) Config {
	return Config{
		MaxAnalyzeAPICall:       maxAnalyzeAPICall,
		BaseURL:                 baseURL,
		OpenAIToken:             openaiToken,
		ChatGPTModel:            chatGPTModel,
		ReasoningModel:          reasoningModel,
		defaultLocale:           defaultLocale,
		SlackAPIToken:           SlackAPIToken,
		excludeDeleteDataSource: excludeDeleteDataSource,
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
	oimac, err := s.newOrganizationIAMClient(clientAddr)
	if err != nil {
		return err
	}
	oc, err := s.newOrganizationClient(clientAddr)
	if err != nil {
		return err
	}
	isvc := iamserver.NewIAMService(s.db, fc, oc, oimac, s.logger)
	asvc := alertserver.NewAlertService(
		s.config.MaxAnalyzeAPICall,
		s.config.BaseURL,
		fc,
		pc,
		iamc,
		s.db,
		s.logger,
		s.config.defaultLocale,
		s.config.SlackAPIToken,
	)
	oisvc := organization_iamserver.NewOrganizationIAMService(s.db, iamc, s.logger)
	fsvc := findingserver.NewFindingService(s.db, s.config.OpenAIToken, s.config.ChatGPTModel, s.config.ReasoningModel, s.config.excludeDeleteDataSource, s.logger)
	psvc := projectserver.NewProjectService(s.db, iamc, s.logger)
	rsvc := reportserver.NewReportService(s.db, s.logger)
	aisvc := aiserver.NewAIService(s.db, s.config.OpenAIToken, s.config.ChatGPTModel, s.config.ReasoningModel, s.logger)
	osvc := organizationserver.NewOrganizationService(s.db, oimac, s.logger)
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
	ai.RegisterAIServiceServer(server, aisvc)
	organization.RegisterOrganizationServiceServer(server, osvc)
	organization_iam.RegisterOrganizationIAMServiceServer(server, oisvc)
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

func (s *Server) newOrganizationIAMClient(svcAddr string) (organization_iam.OrganizationIAMServiceClient, error) {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc connection: err=%w", err)
	}
	return organization_iam.NewOrganizationIAMServiceClient(conn), nil
}

func (s *Server) newOrganizationClient(svcAddr string) (organization.OrganizationServiceClient, error) {
	ctx := context.Background()
	conn, err := getGRPCConn(ctx, svcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc connection: err=%w", err)
	}
	return organization.NewOrganizationServiceClient(conn), nil
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
