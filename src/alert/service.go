package main

import (
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/project"
)

type alertService struct {
	repository           alertRepository
	findingClient        finding.FindingServiceClient
	projectClient        project.ProjectServiceClient
	maxAnalyzeAPICall    int64
	notificationAlertURL string
}
