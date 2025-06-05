package organization_iam

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/organization_iam"
	"gorm.io/gorm"
)

var (
	actionNamePattern = regexp.MustCompile(`^(\w|-)+/(\w|-)+$`)
)

func (i *OrganizationIAMService) IsAuthorizedOrganization(ctx context.Context, req *organization_iam.IsAuthorizedOrganizationRequest) (*organization_iam.IsAuthorizedOrganizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if !actionNamePattern.MatchString(req.ActionName) {
		return nil, fmt.Errorf("invalid action name, pattern=%s, action_name=%s", actionNamePattern, req.ActionName)
	}
	policies, err := i.repository.GetOrganizationPolicyByUserID(ctx, req.UserId, req.OrganizationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.IsAuthorizedOrganizationResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized, err := isAuthorizedByOrganizationPolicy(req.ActionName, policies)
	if err != nil {
		return &organization_iam.IsAuthorizedOrganizationResponse{Ok: false}, err
	}
	if isAuthorized {
		i.logger.Infof(ctx, "Authorized organization action, request=%+v", req)
	}
	return &organization_iam.IsAuthorizedOrganizationResponse{Ok: isAuthorized}, nil
}

func isAuthorizedByOrganizationPolicy(action string, policies *[]model.OrganizationPolicy) (bool, error) {
	for _, p := range *policies {
		actionPtn, err := regexp.Compile(p.ActionPtn)
		if err != nil {
			return false, err
		}
		if actionPtn.MatchString(action) {
			return true, nil
		}
	}
	return false, nil
}

func (o *OrganizationIAMService) IsAdmin(ctx context.Context, req *organization_iam.IsAdminRequest) (*organization_iam.IsAdminResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	policy, err := o.repository.GetAdminOrganizationPolicy(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.IsAdminResponse{Ok: false}, nil
		}
		return nil, err
	}
	if policy == nil || len(*policy) < 1 {
		return &organization_iam.IsAdminResponse{Ok: false}, nil
	}

	// Check if user has organization-admin policy in Admin Organization (ID=1)
	for _, p := range *policy {
		if p.OrganizationID == 1 && p.Name == "organization-admin" {
			o.logger.Debugf(ctx, "user(%d) is admin with organization-admin policy in Admin Organization", req.UserId)
			return &organization_iam.IsAdminResponse{Ok: true}, nil
		}
	}

	o.logger.Debugf(ctx, "user(%d) has Admin Organization policies but no organization-admin role, policies: %d", req.UserId, len(*policy))
	return &organization_iam.IsAdminResponse{Ok: false}, nil
}
