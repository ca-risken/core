package main

import (
	"fmt"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/common/pkg/profiler"
	"github.com/ca-risken/common/pkg/tracer"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/server"
	"github.com/gassara-kys/envconfig"
)

const (
	nameSpace   = "core"
	serviceName = "core"
)

type AppConf struct {
	Port            string   `default:"8080"`
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	TraceDebug      bool     `split_words:"true" default:"false"`

	// service
	MaxAnalyzeAPICall    int64  `split_words:"true" default:"10"`
	NotificationAlertURL string `split_words:"true" default:"http://localhost"`

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
}

func main() {
	var logger = logging.NewLogger()
	var conf AppConf
	err := envconfig.Process("", &conf)
	if err != nil {
		logger.Fatal(err.Error())
	}

	pTypes, err := profiler.ConvertProfileTypeFrom(conf.ProfileTypes)
	if err != nil {
		logger.Fatal(err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(conf.ProfileExporter)
	if err != nil {
		logger.Fatal(err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      conf.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer pc.Stop()

	tc := &tracer.Config{
		ServiceName: getFullServiceName(),
		Environment: conf.EnvName,
		Debug:       conf.TraceDebug,
	}
	tracer.Start(tc)
	defer tracer.Stop()

	dbConf := &db.Config{
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
	db := db.NewClient(dbConf, logger)
	c := server.NewConfig(conf.MaxAnalyzeAPICall, conf.NotificationAlertURL)
	server := server.NewServer("0.0.0.0", conf.Port, db, logger, c)

	err = server.Run()
	if err != nil {
		logger.Fatalf("failed to run server: %w", err)
	}
}

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}
