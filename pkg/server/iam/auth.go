package iam

import (
	"context"
	"errors"
	"regexp"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func (i *IAMService) IsAuthorized(ctx context.Context, req *iam.IsAuthorizedRequest) (*iam.IsAuthorizedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	isAdmin, err := i.isUserAdmin(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		i.logger.Infof(ctx, "Authorized admin user action, request=%+v", req)
		return &iam.IsAuthorizedResponse{Ok: true}, nil
	}

	// Check project-level policies
	policies, err := i.repository.GetUserPolicy(ctx, req.UserId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var isAuthorizedByProject bool
	if policies != nil {
		isAuthorizedByProject, err = isAuthorizedByPolicy(req.ProjectId, req.ActionName, req.ResourceName, policies)
		if err != nil {
			return &iam.IsAuthorizedResponse{Ok: false}, err
		}
		if isAuthorizedByProject {
			i.logger.Infof(ctx, "Authorized user action by project policy, request=%+v", req)
			return &iam.IsAuthorizedResponse{Ok: true}, nil
		}
	}

	// Check organization-level policies for the project
	isAuthorizedByOrg, err := i.checkOrganizationAuthorization(ctx, req.UserId, req.ProjectId, req.ActionName)
	if err != nil {
		// Log error but don't fail completely if organization check fails
		i.logger.Warnf(ctx, "Organization authorization check failed: %v", err)
		isAuthorizedByOrg = false
	}
	if isAuthorizedByOrg {
		i.logger.Infof(ctx, "Authorized user action by organization policy, request=%+v", req)
		return &iam.IsAuthorizedResponse{Ok: true}, nil
	}

	i.logger.Debugf(ctx, "User not authorized: user_id=%d, project_id=%d, action=%s, resource=%s", req.UserId, req.ProjectId, req.ActionName, req.ResourceName)
	return &iam.IsAuthorizedResponse{Ok: false}, nil
}

func (i *IAMService) IsAuthorizedAdmin(ctx context.Context, req *iam.IsAuthorizedAdminRequest) (*iam.IsAuthorizedAdminResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	isAdmin, err := i.isUserAdmin(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	i.logger.Infof(ctx, "Authorized user action, request=%+v", req)
	return &iam.IsAuthorizedAdminResponse{Ok: isAdmin}, nil
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
	isAdmin, err := i.isUserAdmin(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	i.logger.Debugf(ctx, "user(%d) is_admin: %t", req.UserId, isAdmin)
	return &iam.IsAdminResponse{Ok: isAdmin}, nil
}

func (i *IAMService) isUserAdmin(ctx context.Context, userID uint32) (bool, error) {
	user, err := i.repository.GetUser(ctx, userID, "", "")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return user.IsAdmin, nil
}

func (i *IAMService) checkOrganizationAuthorization(ctx context.Context, userID, projectID uint32, actionName string) (bool, error) {
	if i.organizationClient == nil || i.organizationIamClient == nil {
		i.logger.Debugf(ctx, "Organization clients not available, skipping organization authorization check")
		return false, nil
	}

	i.logger.Debugf(ctx, "Checking organization authorization: user_id=%d, project_id=%d, action=%s", userID, projectID, actionName)

	orgList, err := i.organizationClient.ListOrganization(ctx, &organization.ListOrganizationRequest{
		ProjectId: projectID,
	})
	if err != nil {
		i.logger.Warnf(ctx, "Failed to list organizations for project %d: %v", projectID, err)
		return false, err
	}

	i.logger.Debugf(ctx, "Found %d organizations for project %d", len(orgList.Organization), projectID)

	for _, org := range orgList.Organization {
		i.logger.Debugf(ctx, "Checking authorization for organization %d (%s)", org.OrganizationId, org.Name)

		isAuthorized, err := i.organizationIamClient.IsAuthorizedOrganization(ctx, &organization_iam.IsAuthorizedOrganizationRequest{
			UserId:         userID,
			OrganizationId: org.OrganizationId,
			ActionName:     actionName,
		})
		if err != nil {
			i.logger.Warnf(ctx, "Failed to check organization authorization: org_id=%d, user_id=%d, action=%s, error=%v", org.OrganizationId, userID, actionName, err)
			continue
		}

		i.logger.Debugf(ctx, "Organization %d authorization result: %t", org.OrganizationId, isAuthorized.Ok)

		if isAuthorized.Ok {
			i.logger.Infof(ctx, "User authorized through organization: user_id=%d, organization_id=%d, action=%s", userID, org.OrganizationId, actionName)
			return true, nil
		}
	}

	i.logger.Debugf(ctx, "User not authorized through any organization: user_id=%d, project_id=%d, action=%s", userID, projectID, actionName)
	return false, nil
}
