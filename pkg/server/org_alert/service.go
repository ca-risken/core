package org_alert

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/org_alert"
	"github.com/slack-go/slack"
)

var _ org_alert.OrgAlertServiceServer = (*OrgAlertService)(nil)

type OrgAlertService struct {
	repository    db.OrgAlertRepository
	logger        logging.Logger
	slackClient   slack.Client
	defaultLocale string
}

func NewOrgAlertService(
	repository db.OrgAlertRepository,
	logger logging.Logger,
	slackApiToken string,
	defaultLocale string,
) *OrgAlertService {
	return &OrgAlertService{
		repository:    repository,
		logger:        logger,
		slackClient:   *slack.New(slackApiToken),
		defaultLocale: defaultLocale,
	}
}
