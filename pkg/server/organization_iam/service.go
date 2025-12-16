package organization_iam

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/organization_iam"
)

var _ organization_iam.OrganizationIAMServiceServer = (*OrganizationIAMService)(nil)

type OrganizationIAMService struct {
	repository db.OrganizationIAMRepository
	orgClient  organization.OrganizationServiceClient
	iamClient  iam.IAMServiceClient
	logger     logging.Logger
}

func NewOrganizationIAMService(repository db.OrganizationIAMRepository, orgClient organization.OrganizationServiceClient, iamClient iam.IAMServiceClient, logger logging.Logger) *OrganizationIAMService {
	return &OrganizationIAMService{
		repository: repository,
		orgClient:  orgClient,
		iamClient:  iamClient,
		logger:     logger,
	}
}
