package org_alert

import (
	"context"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/org_alert"
	"github.com/cenkalti/backoff/v4"
	"github.com/slack-go/slack"
)

var _ org_alert.OrgAlertServiceServer = (*OrgAlertService)(nil)

type OrgAlertService struct {
	repository    db.OrgAlertRepository
	logger        logging.Logger
	slackClient   slack.Client
	defaultLocale string
	retryer       backoff.BackOff
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
		retryer:       backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10),
	}
}

func (s *OrgAlertService) newRetryLogger(ctx context.Context, funcName string) func(error, time.Duration) {
	return func(err error, t time.Duration) {
		s.logger.Warnf(ctx, "[RetryLogger] %s error: duration=%+v, err=%+v", funcName, t, err)
	}
}
