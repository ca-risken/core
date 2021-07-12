package main

import (
	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/CyberAgent/mimosa-core/proto/project"
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
