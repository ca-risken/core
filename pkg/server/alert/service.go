package alert

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/project"
)

type AlertService struct {
	repository        db.AlertRepository
	findingClient     finding.FindingServiceClient
	projectClient     project.ProjectServiceClient
	maxAnalyzeAPICall int64
	baseURL           string
	logger            logging.Logger
}

func NewAlertService(
	maxAnalyzeAPICall int64,
	baseURL string,
	findingClient finding.FindingServiceClient,
	projectClient project.ProjectServiceClient,
	repository db.AlertRepository,
	logger logging.Logger,
) *AlertService {
	return &AlertService{
		repository:        repository,
		findingClient:     findingClient,
		projectClient:     projectClient,
		maxAnalyzeAPICall: maxAnalyzeAPICall,
		baseURL:           baseURL,
		logger:            logger,
	}
}
