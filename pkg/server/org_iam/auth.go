package org_iam

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/org_iam"
	"gorm.io/gorm"
)

var (
	actionNamePattern = regexp.MustCompile(`^(\w|-)+/(\w|-)+$`)
)

func (i *OrgIAMService) IsAuthorizedOrg(ctx context.Context, req *org_iam.IsAuthorizedOrgRequest) (*org_iam.IsAuthorizedOrgResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	isAdmin, err := i.iamClient.IsAdmin(ctx, &iam.IsAdminRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	if isAdmin.Ok {
		i.logger.Infof(ctx, "Authorized admin user action, request=%+v", req)
		return &org_iam.IsAuthorizedOrgResponse{Ok: true}, nil
	}
	if !actionNamePattern.MatchString(req.ActionName) {
		return nil, fmt.Errorf("invalid action name, pattern=%s, action_name=%s", actionNamePattern, req.ActionName)
	}
	policies, err := i.repository.GetOrgPolicyByUserID(ctx, req.UserId, req.OrganizationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_iam.IsAuthorizedOrgResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized, err := isAuthorizedByOrgPolicy(req.ActionName, req.ProjectName, policies)
	if err != nil {
		return &org_iam.IsAuthorizedOrgResponse{Ok: false}, err
	}
	if isAuthorized {
		i.logger.Infof(ctx, "Authorized organization action, request=%+v", req)
	}
	return &org_iam.IsAuthorizedOrgResponse{Ok: isAuthorized}, nil
}

func (i *OrgIAMService) IsAuthorizedOrgToken(ctx context.Context, req *org_iam.IsAuthorizedOrgTokenRequest) (*org_iam.IsAuthorizedOrgTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if !actionNamePattern.MatchString(req.ActionName) {
		return nil, fmt.Errorf("invalid action name, pattern=%s, action_name=%s", actionNamePattern, req.ActionName)
	}
	if req.ProjectId != 0 {
		orgList, err := i.orgClient.ListOrganization(ctx, &organization.ListOrganizationRequest{
			ProjectId: req.ProjectId,
		})
		if err != nil {
			i.logger.Warnf(ctx, "Failed to list organizations for project %d: %v", req.ProjectId, err)
			return &org_iam.IsAuthorizedOrgTokenResponse{Ok: false}, nil
		}
		isExists := false
		for _, org := range orgList.Organization {
			if org.OrganizationId == req.OrganizationId {
				isExists = true
				break
			}
		}
		if !isExists {
			return &org_iam.IsAuthorizedOrgTokenResponse{Ok: false}, nil
		}
	}
	existsMaintainer, err := i.repository.ExistsOrgAccessTokenMaintainer(ctx, req.OrganizationId, req.AccessTokenId)
	if err != nil {
		return nil, err
	}
	if !existsMaintainer {
		i.logger.Warnf(ctx, "Unauthorized organization access token that has no maintainers or expired. organization_id=%d, access_token_id=%d", req.OrganizationId, req.AccessTokenId)
		return &org_iam.IsAuthorizedOrgTokenResponse{Ok: false}, nil
	}
	policies, err := i.repository.GetOrgTokenPolicy(ctx, req.OrganizationId, req.AccessTokenId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_iam.IsAuthorizedOrgTokenResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized, err := isAuthorizedByOrgPolicy(req.ActionName, req.ProjectName, policies)
	if err != nil {
		return &org_iam.IsAuthorizedOrgTokenResponse{Ok: false}, err
	}
	if isAuthorized {
		i.logger.Infof(ctx, "Authorized organization access token action, request=%+v", req)
	}
	return &org_iam.IsAuthorizedOrgTokenResponse{Ok: isAuthorized}, nil
}

func isAuthorizedByOrgPolicy(action string, projectName string, policies *[]model.OrganizationPolicy) (bool, error) {
	for _, p := range *policies {
		actionPtn, err := regexp.Compile(p.ActionPtn)
		if err != nil {
			return false, err
		}
		projectPtn, err := regexp.Compile(p.ProjectPtn)
		if err != nil {
			return false, err
		}
		if actionPtn.MatchString(action) && projectPtn.MatchString(projectName) {
			return true, nil
		}
	}
	return false, nil
}
