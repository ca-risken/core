package main

import (
	"fmt"
	"net"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/common/pkg/profiler"
	mimosarpc "github.com/ca-risken/common/pkg/rpc"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/core/proto/iam"
	"github.com/gassara-kys/envconfig"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

const (
	nameSpace   = "core"
	serviceName = "iam"
)

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}

type AppConf struct {
	Debug           string   `default:"false"`
	Port            string   `default:"8002"`
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

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

	// grpc
	FindingSvcAddr string `required:"true" split_words:"true" default:"finding.core.svc.cluster.local:8001"`
}

func main() {
	var conf AppConf
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatal(err.Error())
	}
	if conf.Debug == "true" {
		appLogger.Level(logging.DebugLevel)
	}
	appLogger.Infof("Load IAM config: %+v", conf)

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

	tc := &tracer.Config{
		ServiceName: getFullServiceName(),
		Environment: conf.EnvName,
		Debug:       conf.TraceDebug,
	}
	tracer.Start(tc)
	defer tracer.Stop()

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		appLogger.Fatal(err)
	}

	service := &iamService{}
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
	service.repository = newIAMRepository(dbConf)
	service.findingClient = newFindingClient(conf.FindingSvcAddr)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				mimosarpc.LoggingUnaryServerInterceptor(appLogger),
				grpctrace.UnaryServerInterceptor())))
	iam.RegisterIAMServiceServer(server, service)

	reflection.Register(server) // enable reflection API
	appLogger.Infof("Starting gRPC server at :%s", conf.Port)
	if err := server.Serve(l); err != nil {
		appLogger.Fatalf("Failed to gRPC serve: %v", err)
	}
}
