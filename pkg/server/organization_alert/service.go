package organization_alert

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/organization_alert"
	"github.com/slack-go/slack"
)

var _ organization_alert.OrganizationAlertServiceServer = (*OrganizationAlertService)(nil)

type OrganizationAlertService struct {
	repository    db.OrganizationAlertRepository
	logger        logging.Logger
	slackClient   slack.Client
	defaultLocale string
	baseURL       string
}

func NewOrganizationAlertService(
	repository db.OrganizationAlertRepository,
	logger logging.Logger,
	defaultLocale string,
	baseURL string,
	slackApiToken string,
) *OrganizationAlertService {
	return &OrganizationAlertService{
		repository:    repository,
		logger:        logger,
		slackClient:   *slack.New(slackApiToken),
		defaultLocale: defaultLocale,
		baseURL:       baseURL,
	}
}
