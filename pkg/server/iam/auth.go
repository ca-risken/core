package iam

import (
	"context"
	"errors"
	"regexp"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func (i *IAMService) IsAuthorized(ctx context.Context, req *iam.IsAuthorizedRequest) (*iam.IsAuthorizedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	policies, err := i.repository.GetUserPolicy(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &iam.IsAuthorizedResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized, err := isAuthorizedByPolicy(req.ProjectId, req.ActionName, req.ResourceName, policies)
	if err != nil {
		return &iam.IsAuthorizedResponse{Ok: false}, err
	}
	if isAuthorized {
		i.logger.Infof(ctx, "Authorized user action, request=%+v", req)
	}
	return &iam.IsAuthorizedResponse{Ok: isAuthorized}, nil
}

func (i *IAMService) IsAuthorizedAdmin(ctx context.Context, req *iam.IsAuthorizedAdminRequest) (*iam.IsAuthorizedAdminResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check Project ID=null IAM policies
	policies, err := i.repository.GetAdminPolicy(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			i.logger.Debugf(ctx, "user(%d) has no admin IAM policies", req.UserId)
			return &iam.IsAuthorizedAdminResponse{Ok: false}, nil
		}
		return nil, err
	}
	if policies == nil || len(*policies) < 1 {
		i.logger.Debugf(ctx, "user(%d) has no admin IAM policies", req.UserId)
		return &iam.IsAuthorizedAdminResponse{Ok: false}, nil
	}

	// Check Organization ID=null Organization IAM policies
	orgIamReq := &organization_iam.GetSystemAdminOrganizationPolicyRequest{
		UserId: req.UserId,
	}
	orgIamResp, err := i.organizationIamClient.GetSystemAdminOrganizationPolicy(ctx, orgIamReq)
	if err != nil {
		i.logger.Errorf(ctx, "failed to get system admin organization policies for user(%d): %v", req.UserId, err)
		return nil, err
	}
	if orgIamResp == nil || len(orgIamResp.OrganizationPolicies) < 1 {
		i.logger.Debugf(ctx, "user(%d) has no system admin organization IAM policies", req.UserId)
		return &iam.IsAuthorizedAdminResponse{Ok: false}, nil
	}

	// User has both Project ID=null IAM policies AND Organization ID=null Organization IAM policies
	// Now check if authorized for the specific action and resource
	isAuthorized, err := isAuthorizedByPolicy(0, req.ActionName, req.ResourceName, policies)
	if err != nil {
		return &iam.IsAuthorizedAdminResponse{Ok: false}, err
	}
	if isAuthorized {
		i.logger.Infof(ctx, "Authorized admin action, request=%+v", req)
	}
	return &iam.IsAuthorizedAdminResponse{Ok: isAuthorized}, nil
}

func (i *IAMService) IsAuthorizedToken(ctx context.Context, req *iam.IsAuthorizedTokenRequest) (*iam.IsAuthorizedTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	existsMaintainer, err := i.repository.ExistsAccessTokenMaintainer(ctx, req.ProjectId, req.AccessTokenId)
	if err != nil {
		return nil, err
	}
	if !existsMaintainer {
		i.logger.Warnf(ctx, "Unautorized the token that has no maintainers or expired in the project. project_id=%d, access_token_id=%d", req.ProjectId, req.AccessTokenId)
		return &iam.IsAuthorizedTokenResponse{Ok: false}, nil
	}
	policies, err := i.repository.GetTokenPolicy(ctx, req.AccessTokenId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &iam.IsAuthorizedTokenResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized, err := isAuthorizedByPolicy(req.ProjectId, req.ActionName, req.ResourceName, policies)
	if err != nil {
		return &iam.IsAuthorizedTokenResponse{Ok: false}, err
	}
	if isAuthorized {
		i.logger.Infof(ctx, "Authorized access_token action, request=%+v", req)
	}
	return &iam.IsAuthorizedTokenResponse{Ok: isAuthorized}, nil
}

func isAuthorizedByPolicy(projectID uint32, action, resource string, policies *[]model.Policy) (bool, error) {
	for _, p := range *policies {
		actionPtn, err := regexp.Compile(p.ActionPtn)
		if err != nil {
			return false, err
		}
		resourcePtn, err := regexp.Compile(p.ResourcePtn)
		if err != nil {
			return false, err
		}
		if !zero.IsZeroVal(p.ProjectID) && projectID != p.ProjectID {
			continue
		}
		if actionPtn.MatchString(action) && resourcePtn.MatchString(resource) {
			return true, nil
		}
	}
	return false, nil

}

func (i *IAMService) IsAdmin(ctx context.Context, req *iam.IsAdminRequest) (*iam.IsAdminResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	iamPolicy, err := i.repository.GetAdminPolicy(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			i.logger.Debugf(ctx, "user(%d) has no admin IAM policies", req.UserId)
			return &iam.IsAdminResponse{Ok: false}, nil
		}
		return nil, err
	}
	if iamPolicy == nil || len(*iamPolicy) < 1 {
		i.logger.Debugf(ctx, "user(%d) has no admin IAM policies", req.UserId)
		return &iam.IsAdminResponse{Ok: false}, nil
	}
	orgIamReq := &organization_iam.GetSystemAdminOrganizationPolicyRequest{
		UserId: req.UserId,
	}
	orgIamResp, err := i.organizationIamClient.GetSystemAdminOrganizationPolicy(ctx, orgIamReq)
	if err != nil {
		i.logger.Errorf(ctx, "failed to get system admin organization policies for user(%d): %v", req.UserId, err)
		return nil, err
	}
	if orgIamResp == nil || len(orgIamResp.OrganizationPolicies) < 1 {
		i.logger.Debugf(ctx, "user(%d) has no system admin organization IAM policies", req.UserId)
		return &iam.IsAdminResponse{Ok: false}, nil
	}
	i.logger.Infof(ctx, "user(%d) is system admin - has %d IAM policies and %d organization IAM policies",
		req.UserId, len(*iamPolicy), len(orgIamResp.OrganizationPolicies))
	return &iam.IsAdminResponse{Ok: true}, nil
}
