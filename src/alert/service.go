package main

import (
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/project"
)

type alertService struct {
	repository    alertRepository
	findingClient finding.FindingServiceClient
	projectClient project.ProjectServiceClient
}

func newAlertService() alert.AlertServiceServer {
	return &alertService{
		repository:    newAlertRepository(),
		findingClient: newFindingClient(),
		projectClient: newProjectClient(),
	}
}
