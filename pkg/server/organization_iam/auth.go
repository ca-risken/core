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

func (i *OrganizationIAMService) GetSystemAdminOrganizationPolicy(ctx context.Context, req *organization_iam.GetSystemAdminOrganizationPolicyRequest) (*organization_iam.GetSystemAdminOrganizationPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	policies, err := i.repository.GetSystemAdminOrganizationPolicy(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.GetSystemAdminOrganizationPolicyResponse{}, nil
		}
		return nil, err
	}
	var orgPolicies []*organization_iam.OrganizationPolicy
	for _, p := range *policies {
		orgPolicies = append(orgPolicies, &organization_iam.OrganizationPolicy{
			PolicyId:       p.PolicyID,
			Name:           p.Name,
			OrganizationId: p.OrganizationID,
			ActionPtn:      p.ActionPtn,
			CreatedAt:      p.CreatedAt.Unix(),
			UpdatedAt:      p.UpdatedAt.Unix(),
		})
	}
	return &organization_iam.GetSystemAdminOrganizationPolicyResponse{OrganizationPolicies: orgPolicies}, nil
}
