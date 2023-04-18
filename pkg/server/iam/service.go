package iam

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/iam"
)

var _ iam.IAMServiceServer = (*IAMService)(nil)

type IAMService struct {
	repository    db.IAMRepository
	findingClient finding.FindingServiceClient
	logger        logging.Logger
}

func NewIAMService(repository db.IAMRepository, findingClient finding.FindingServiceClient, logger logging.Logger) *IAMService {
	return &IAMService{
		repository:    repository,
		findingClient: findingClient,
		logger:        logger,
	}
}
