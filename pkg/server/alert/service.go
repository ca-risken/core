package alert

import (
	"context"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/alert"
	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/org_alert"
	"github.com/ca-risken/core/proto/project"
	"github.com/cenkalti/backoff/v4"
	"github.com/slack-go/slack"
)

var _ alert.AlertServiceServer = (*AlertService)(nil)

type AlertService struct {
	repository        db.AlertRepository
	findingClient     finding.FindingServiceClient
	projectClient     project.ProjectServiceClient
	iamClient         iam.IAMServiceClient
	orgAlertClient    org_alert.OrgAlertServiceClient
	maxAnalyzeAPICall int64
	baseURL           string
	logger            logging.Logger
	defaultLocale     string
	aiSummaryEnabled  bool
	summaryLanguage   string
	slackClient       slack.Client
	retryer           backoff.BackOff
}

func NewAlertService(
	maxAnalyzeAPICall int64,
	baseURL string,
	findingClient finding.FindingServiceClient,
	projectClient project.ProjectServiceClient,
	iamClient iam.IAMServiceClient,
	orgAlertClient org_alert.OrgAlertServiceClient,
	repository db.AlertRepository,
	logger logging.Logger,
	defaultLocale string,
	aiSummaryEnabled bool,
	summaryLanguage string,
	slackApiToken string,
) *AlertService {
	return &AlertService{
		repository:        repository,
		findingClient:     findingClient,
		projectClient:     projectClient,
		iamClient:         iamClient,
		orgAlertClient:    orgAlertClient,
		maxAnalyzeAPICall: maxAnalyzeAPICall,
		baseURL:           baseURL,
		logger:            logger,
		defaultLocale:     defaultLocale,
		aiSummaryEnabled:  aiSummaryEnabled,
		summaryLanguage:   summaryLanguage,
		slackClient:       *slack.New(slackApiToken),
		retryer:           backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 10),
	}
}

func (a *AlertService) newRetryLogger(ctx context.Context, funcName string) func(error, time.Duration) {
	return func(err error, t time.Duration) {
		a.logger.Warnf(ctx, "[RetryLogger] %s error: duration=%+v, err=%+v", funcName, t, err)
	}
}
