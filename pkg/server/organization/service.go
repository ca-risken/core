package organization

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/org_iam"
)

type OrganizationService struct {
	repository            db.OrganizationRepository
	orgIamClient org_iam.OrgIAMServiceClient
	logger                logging.Logger
}

func NewOrganizationService(repository db.OrganizationRepository, orgIamClient org_iam.OrgIAMServiceClient, logger logging.Logger) *OrganizationService {
	return &OrganizationService{
		repository:            repository,
		orgIamClient: orgIamClient,
		logger:                logger,
	}
}
