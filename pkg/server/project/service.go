package project

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/project"
)

var _ project.ProjectServiceServer = (*ProjectService)(nil)

type ProjectService struct {
	repository db.ProjectRepository
	iamClient  iam.IAMServiceClient
	logger     logging.Logger
}

func NewProjectService(repository db.ProjectRepository, iamClient iam.IAMServiceClient, logger logging.Logger) *ProjectService {
	return &ProjectService{
		repository: repository,
		iamClient:  iamClient,
		logger:     logger,
	}
}
