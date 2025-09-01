package main

import (
	"context"
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

var (
	samplingRate float64 = 0.3000
)

type AppConf struct {
	Port            string   `default:"8080"`
	EnvName         string   `default:"local" split_words:"true"`
	ProfileExporter string   `split_words:"true" default:"nop"`
	ProfileTypes    []string `split_words:"true"`
	Debug           bool     `split_words:"true" default:"false"`
	TraceDebug      bool     `split_words:"true" default:"false"`

	// service
	MaxAnalyzeAPICall       int64    `split_words:"true" default:"10"`
	BaseURL                 string   `split_words:"true" default:"http://localhost"`
	OpenAIToken             string   `split_words:"true"`
	ChatGPTModel            string   `split_words:"true" default:"gpt-4.1"`
	ReasoningModel          string   `split_words:"true" default:"gpt-5"`
	DefaultLocale           string   `split_words:"true" default:"en"`
	SlackAPIToken           string   `split_words:"true"`
	ExcludeDeleteDataSource []string `split_words:"true" default:"code:gitleaks"`

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
	ctx := context.Background()
	var logger = logging.NewLogger()
	var conf AppConf
	err := envconfig.Process("", &conf)
	if err != nil {
		logger.Fatal(ctx, err.Error())
	}
	if conf.Debug {
		logger.Level(logging.DebugLevel)
	}

	pTypes, err := profiler.ConvertProfileTypeFrom(conf.ProfileTypes)
	if err != nil {
		logger.Fatal(ctx, err.Error())
	}
	pExporter, err := profiler.ConvertExporterTypeFrom(conf.ProfileExporter)
	if err != nil {
		logger.Fatal(ctx, err.Error())
	}
	pc := profiler.Config{
		ServiceName:  getFullServiceName(),
		EnvName:      conf.EnvName,
		ProfileTypes: pTypes,
		ExporterType: pExporter,
	}
	err = pc.Start()
	if err != nil {
		logger.Fatal(ctx, err.Error())
	}
	defer pc.Stop()

	tc := &tracer.Config{
		ServiceName:  getFullServiceName(),
		Environment:  conf.EnvName,
		Debug:        conf.TraceDebug,
		SamplingRate: &samplingRate,
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
	db, err := db.NewClient(dbConf, logger)
	if err != nil {
		logger.Fatalf(ctx, "failed to create database client: %w", err)
	}
	c := server.NewConfig(conf.MaxAnalyzeAPICall, conf.BaseURL, conf.OpenAIToken, conf.ChatGPTModel, conf.ReasoningModel, conf.DefaultLocale, conf.SlackAPIToken, conf.ExcludeDeleteDataSource)
	server := server.NewServer("0.0.0.0", conf.Port, db, logger, c)

	err = server.Run(ctx)
	if err != nil {
		logger.Fatalf(ctx, "failed to run server: %w", err)
	}
}

func getFullServiceName() string {
	return fmt.Sprintf("%s.%s", nameSpace, serviceName)
}
