package organization

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
)

type OrganizationService struct {
	repository db.OrganizationRepository
	logger     logging.Logger
}

func NewOrganizationService(repository db.OrganizationRepository, logger logging.Logger) *OrganizationService {
	return &OrganizationService{
		repository: repository,
		logger:     logger,
	}
}
