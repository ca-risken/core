package org_alert

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/org_alert"
)

var _ org_alert.OrgAlertServiceServer = (*OrgAlertService)(nil)

type OrgAlertService struct {
	repository db.OrgAlertRepository
	logger     logging.Logger
}

func NewOrgAlertService(
	repository db.OrgAlertRepository,
	logger logging.Logger,
) *OrgAlertService {
	return &OrgAlertService{
		repository: repository,
		logger:     logger,
	}
}
