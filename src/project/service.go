package main

import (
	"github.com/ca-risken/core/proto/project"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
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
		appLogger.SetLevel(logrus.DebugLevel)
	}
	return &projectService{
		repository: newProjectRepository(),
		iamClient:  newIAMService(),
	}
}
