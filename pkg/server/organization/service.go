package organization

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/organization_iam"
)

type OrganizationService struct {
	repository            db.OrganizationRepository
	organizationIamClient organization_iam.OrganizationIAMServiceClient
	logger                logging.Logger
}

func NewOrganizationService(repository db.OrganizationRepository, organizationIamClient organization_iam.OrganizationIAMServiceClient, logger logging.Logger) *OrganizationService {
	return &OrganizationService{
		repository:            repository,
		organizationIamClient: organizationIamClient,
		logger:                logger,
	}
}
