package project

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/ca-risken/core/proto/project"
)

var _ project.ProjectServiceServer = (*ProjectService)(nil)

type ProjectService struct {
	repository            db.ProjectRepository
	iamClient             iam.IAMServiceClient
	organizationClient    organization.OrganizationServiceClient
	organizationIamClient organization_iam.OrganizationIAMServiceClient
	logger                logging.Logger
}

func NewProjectService(repository db.ProjectRepository, iamClient iam.IAMServiceClient, organizationClient organization.OrganizationServiceClient, organizationIamClient organization_iam.OrganizationIAMServiceClient, logger logging.Logger) *ProjectService {
	return &ProjectService{
		repository:            repository,
		iamClient:             iamClient,
		organizationClient:    organizationClient,
		organizationIamClient: organizationIamClient,
		logger:                logger,
	}
}
