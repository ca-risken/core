package alert

import (
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/project"
)

type AlertService struct {
	repository           db.AlertRepository
	findingClient        finding.FindingServiceClient
	projectClient        project.ProjectServiceClient
	maxAnalyzeAPICall    int64
	notificationAlertURL string
}

func NewAlertService(maxAnalyzeAPICall int64, notificationAlertURL string, findingClient finding.FindingServiceClient, projectClient project.ProjectServiceClient, repository db.AlertRepository) *AlertService {
	return &AlertService{
		repository:           repository,
		findingClient:        findingClient,
		projectClient:        projectClient,
		maxAnalyzeAPICall:    maxAnalyzeAPICall,
		notificationAlertURL: notificationAlertURL,
	}
}
