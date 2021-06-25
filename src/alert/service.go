package main

import (
	"github.com/CyberAgent/mimosa-core/proto/alert"
	"github.com/CyberAgent/mimosa-core/proto/finding"
)

type alertService struct {
	repository    alertRepository
	findingClient finding.FindingServiceClient
}

func newAlertService() alert.AlertServiceServer {
	return &alertService{
		repository:    newAlertRepository(),
		findingClient: newFindingClient(),
	}
}
