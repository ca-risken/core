package organization_alert

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/organization_alert"
)

var _ organization_alert.OrganizationAlertServiceServer = (*OrganizationAlertService)(nil)

type OrganizationAlertService struct {
	repository db.OrganizationAlertRepository
	logger     logging.Logger
}

func NewOrganizationAlertService(
	repository db.OrganizationAlertRepository,
	logger logging.Logger,
) *OrganizationAlertService {
	return &OrganizationAlertService{
		repository: repository,
		logger:     logger,
	}
}
