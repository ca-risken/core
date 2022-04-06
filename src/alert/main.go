package main

import (
	"fmt"
	"net"

	"github.com/ca-risken/common/pkg/profiler"
	mimosarpc "github.com/ca-risken/common/pkg/rpc"
	"github.com/ca-risken/core/proto/alert"
	"github.com/gassara-kys/envconfig"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	nameSpace   = "core"
	serviceName = "alert"
)

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}

type AppConf struct {
	Port            string   `default:"8004"`
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

	// grpc
	FindingSvcAddr string `required:"true" split_words:"true" default:"finding.core.svc.cluster.local:8001"`
	ProjectSvcAddr string `required:"true" split_words:"true" default:"project.core.svc.cluster.local:8003"`

	// service
	MaxAnalyzeAPICall int64 `split_words:"true" default:"10"`

	// db
	DBMasterHost     string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	DBMasterUser     string `split_words:"true" default:"hoge"`
	DBMasterPassword string `split_words:"true" default:"moge"`
	DBSlaveHost      string `split_words:"true" default:"db.middleware.svc.cluster.local"`
	DBSlaveUser      string `split_words:"true" default:"hoge"`
	DBSlavePassword  string `split_words:"true" default:"moge"`

	DBSchema        string `required:"true"    default:"mimosa"`
	DBPort          int    `required:"true"    default:"3306"`
	DBLogMode       bool   `split_words:"true" default:"false"`
	DBMaxConnection int    `split_words:"true" default:"10"`

	// notification
	NotificationAlertURL string `split_words:"true" default:"http://localhost"`
}

func main() {
	var conf AppConf
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}

	pTypes, err := profiler.ConvertProfileTypeFrom(conf.ProfileTypes)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(conf.ProfileExporter)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      conf.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	defer pc.Stop()

	traceOpts := []tracer.StartOption{
		tracer.WithEnv(conf.EnvName),
		tracer.WithService(getFullServiceName()),
		tracer.WithDebugMode(conf.TraceDebug),
	}
	tracer.Start(traceOpts...)
	defer tracer.Stop()

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		appLogger.Fatal(err)
	}

	service := &alertService{
		maxAnalyzeAPICall:    conf.MaxAnalyzeAPICall,
		notificationAlertURL: conf.NotificationAlertURL,
	}
	service.findingClient = newFindingClient(conf.FindingSvcAddr)
	service.projectClient = newProjectClient(conf.ProjectSvcAddr)
	dbConf := &DBConfig{
		MasterHost:     conf.DBMasterHost,
		MasterUser:     conf.DBMasterUser,
		MasterPassword: conf.DBMasterPassword,
		SlaveHost:      conf.DBSlaveHost,
		SlaveUser:      conf.DBSlaveUser,
		SlavePassword:  conf.DBSlavePassword,
		Schema:         conf.DBSchema,
		Port:           conf.DBPort,
		LogMode:        conf.DBLogMode,
		MaxConnection:  conf.DBMaxConnection,
	}
	service.repository = newAlertRepository(dbConf)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				mimosarpc.LoggingUnaryServerInterceptor(appLogger),
				grpctrace.UnaryServerInterceptor())))
	alert.RegisterAlertServiceServer(server, service)

	reflection.Register(server) // enable reflection API
	appLogger.Infof("Starting gRPC server at :%s", conf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Fatalf("Failed to gRPC serve: %v", err)
	}
}
