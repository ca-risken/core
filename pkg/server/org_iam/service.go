package org_iam

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/org_iam"
)

var _ org_iam.OrgIAMServiceServer = (*OrgIAMService)(nil)

type OrgIAMService struct {
	repository db.OrgIAMRepository
	orgClient  organization.OrganizationServiceClient
	iamClient  iam.IAMServiceClient
	logger     logging.Logger
}

func NewOrgIAMService(repository db.OrgIAMRepository, orgClient organization.OrganizationServiceClient, iamClient iam.IAMServiceClient, logger logging.Logger) *OrgIAMService {
	return &OrgIAMService{
		repository: repository,
		orgClient:  orgClient,
		iamClient:  iamClient,
		logger:     logger,
	}
}
