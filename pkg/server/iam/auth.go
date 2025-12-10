package iam

import (
	"context"
	"errors"
	"regexp"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/organization_iam"
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

	isAuthorizedByProject, err := i.isAuthorizedByProject(ctx, req.UserId, req.ProjectId, req.ActionName, req.ResourceName)
	if err != nil {
		i.logger.Warnf(ctx, "Project authorization check failed: %v", err)
		return nil, err
	}
	if isAuthorizedByProject {
		i.logger.Infof(ctx, "Authorized user action by project policy, request=%+v", req)
		return &iam.IsAuthorizedResponse{Ok: true}, nil
	}

	isAuthorizedByOrg, err := i.isAuthorizedByOrganizations(ctx, req.UserId, req.ProjectId, req.ActionName)
	if err != nil {
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
		if projectID != p.ProjectID {
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

func (i *IAMService) isAuthorizedByOrganizations(ctx context.Context, userID, projectID uint32, actionName string) (bool, error) {
	orgList, err := i.organizationClient.ListOrganization(ctx, &organization.ListOrganizationRequest{
		ProjectId: projectID,
		UserId:    userID,
	})
	if err != nil {
		i.logger.Warnf(ctx, "Failed to list organizations for project %d: %v", projectID, err)
		return false, err
	}
	for _, org := range orgList.Organization {
		isAuthorized, err := i.organizationIamClient.IsAuthorizedOrganization(ctx, &organization_iam.IsAuthorizedOrganizationRequest{
			UserId:         userID,
			OrganizationId: org.OrganizationId,
			ActionName:     actionName,
		})
		if err != nil {
			i.logger.Warnf(ctx, "Failed to check organization authorization: org_id=%d, user_id=%d, action=%s, error=%v", org.OrganizationId, userID, actionName, err)
			continue
		}
		if isAuthorized.Ok {
			i.logger.Infof(ctx, "User authorized through organization: user_id=%d, organization_id=%d, action=%s", userID, org.OrganizationId, actionName)
			return true, nil
		}
	}
	i.logger.Debugf(ctx, "User not authorized through any organization: user_id=%d, project_id=%d, action=%s", userID, projectID, actionName)
	return false, nil
}

func (i *IAMService) isAuthorizedByProject(ctx context.Context, userID, projectID uint32, actionName, resourceName string) (bool, error) {
	policies, err := i.repository.GetUserPolicy(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	if policies != nil {
		isAuthorizedByProject, err := isAuthorizedByPolicy(projectID, actionName, resourceName, policies)
		if err != nil {
			return false, err
		}
		if isAuthorizedByProject {
			return true, nil
		}
	}
	return false, nil
}
