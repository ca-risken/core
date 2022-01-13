package main

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/project"
	"github.com/gassara-kys/envconfig"
)

type projectService struct {
	repository projectRepository
	iamClient  iamService
}

type projectConfig struct {
	Port  string `default:"8003"`
	Debug bool   `default:"false"`
}

func newProjectService() project.ProjectServiceServer {
	var conf projectConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatalf("project config load error: err=%+v", err)
	}
	if conf.Debug {
		appLogger.Level(logging.DebugLevel)
	}
	return &projectService{
		repository: newProjectRepository(),
		iamClient:  newIAMService(),
	}
}
