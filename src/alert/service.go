package main

import (
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/project"
	"github.com/gassara-kys/envconfig"
)

type alertConfig struct {
	MaxAnalyzeAPICall int64 `split_words:"true" default:"10"`
}

type alertService struct {
	maxAnalyzeAPICall int64
	repository        alertRepository
	findingClient     finding.FindingServiceClient
	projectClient     project.ProjectServiceClient
}

func newAlertService() alert.AlertServiceServer {
	var conf alertConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatalf("Faild to load finding config error: err=%+v", err)
	}
	return &alertService{
		maxAnalyzeAPICall: conf.MaxAnalyzeAPICall,
		repository:        newAlertRepository(),
		findingClient:     newFindingClient(),
		projectClient:     newProjectClient(),
	}
}
