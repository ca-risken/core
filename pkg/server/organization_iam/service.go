package organization_iam

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization_iam"
)

var _ organization_iam.OrganizationIAMServiceServer = (*OrganizationIAMService)(nil)

type OrganizationIAMService struct {
	repository db.OrganizationIAMRepository
	iamClient  iam.IAMServiceClient
	logger     logging.Logger
}

func NewOrganizationIAMService(repository db.OrganizationIAMRepository, iamClient iam.IAMServiceClient, logger logging.Logger) *OrganizationIAMService {
	return &OrganizationIAMService{
		repository: repository,
		iamClient:  iamClient,
		logger:     logger,
	}
}
