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
	"github.com/ca-risken/core/proto/project"
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
	projectName, err := i.resolveProjectName(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	policies, err := i.repository.GetOrgPolicyByUserID(ctx, req.UserId, req.OrganizationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_iam.IsAuthorizedOrgResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized, err := isAuthorizedByOrgPolicy(req.ActionName, projectName, policies)
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
	projectName, err := i.resolveProjectName(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	policies, err := i.repository.GetOrgTokenPolicy(ctx, req.OrganizationId, req.AccessTokenId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_iam.IsAuthorizedOrgTokenResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized, err := isAuthorizedByOrgPolicy(req.ActionName, projectName, policies)
	if err != nil {
		return &org_iam.IsAuthorizedOrgTokenResponse{Ok: false}, err
	}
	if isAuthorized {
		i.logger.Infof(ctx, "Authorized organization access token action, request=%+v", req)
	}
	return &org_iam.IsAuthorizedOrgTokenResponse{Ok: isAuthorized}, nil
}

// resolveProjectName resolves project name from project_id via projectClient.
// Returns empty string when projectID is 0 (org-level authorization without project context).
func (i *OrgIAMService) resolveProjectName(ctx context.Context, projectID uint32) (string, error) {
	if projectID == 0 {
		return "", nil
	}
	resp, err := i.projectClient.ListProject(ctx, &project.ListProjectRequest{ProjectId: projectID})
	if err != nil {
		return "", fmt.Errorf("failed to list project, project_id=%d, err=%w", projectID, err)
	}
	for _, p := range resp.Project {
		if p.ProjectId == projectID {
			return p.Name, nil
		}
	}
	return "", fmt.Errorf("project not found, project_id=%d", projectID)
}

// isAuthorizedByOrgPolicy returns true when at least one policy matches the action and project name.
// When projectName is empty (no project context), project_ptn is not evaluated.
func isAuthorizedByOrgPolicy(action, projectName string, policies *[]model.OrganizationPolicy) (bool, error) {
	for _, p := range *policies {
		actionPtn, err := regexp.Compile(p.ActionPtn)
		if err != nil {
			return false, err
		}
		if !actionPtn.MatchString(action) {
			continue
		}
		if projectName == "" {
			return true, nil
		}
		projectPtn, err := regexp.Compile(p.ProjectPtn)
		if err != nil {
			return false, err
		}
		if projectPtn.MatchString(projectName) {
			return true, nil
		}
	}
	return false, nil
}
