package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/ca-risken/core/proto/project"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func convertProjectWithTag(p *db.ProjectWithTag) *project.Project {
	if p == nil {
		return &project.Project{}
	}
	tags := []*project.ProjectTag{}
	if p.Tag != nil {
		for _, t := range *p.Tag {
			tags = append(tags, &project.ProjectTag{
				ProjectId: t.ProjectID,
				Tag:       t.Tag,
				Color:     t.Color,
			})
		}
	}
	return &project.Project{
		ProjectId: p.ProjectID,
		Name:      p.Name,
		Tag:       tags,
		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: p.UpdatedAt.Unix(),
	}
}

func (p *ProjectService) ListProject(ctx context.Context, req *project.ListProjectRequest) (*project.ListProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	directProjects, err := p.repository.ListProject(ctx, req.UserId, req.ProjectId, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	projectMap := make(map[uint32]*db.ProjectWithTag)
	if directProjects != nil {
		for _, pr := range *directProjects {
			projectMap[pr.ProjectID] = &pr
		}
	}
	orgProjects, err := p.getProjectsFromOrganizations(ctx, req.UserId, req.ProjectId, req.Name)
	if err != nil {
		p.logger.Warnf(ctx, "Failed to get projects from organizations: %v", err)
	} else {
		for _, pr := range orgProjects {
			if _, exists := projectMap[pr.ProjectID]; !exists {
				projectMap[pr.ProjectID] = pr
			}
		}
	}
	var prs []*project.Project
	for _, pr := range projectMap {
		prs = append(prs, convertProjectWithTag(pr))
	}
	return &project.ListProjectResponse{Project: prs}, nil
}

// getProjectsFromOrganizations gets projects that the user can access through organization membership
func (p *ProjectService) getProjectsFromOrganizations(ctx context.Context, userID, projectID uint32, name string) ([]*db.ProjectWithTag, error) {
	if p.organizationClient == nil || p.organizationIamClient == nil {
		p.logger.Debugf(ctx, "Organization clients not available, skipping organization project lookup")
		return nil, nil
	}

	// Get organizations where the user has any role
	userOrgs, err := p.getUserOrganizations(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	p.logger.Debugf(ctx, "Found %d organizations for user %d", len(userOrgs), userID)

	var allProjects []*db.ProjectWithTag

	// For each organization, get its projects
	for _, orgID := range userOrgs {
		orgProjectsResp, err := p.organizationClient.ListProjectsInOrganization(ctx, &organization.ListProjectsInOrganizationRequest{
			OrganizationId: orgID,
		})
		if err != nil {
			p.logger.Warnf(ctx, "Failed to get projects for organization %d: %v", orgID, err)
			continue
		}

		// Convert organization projects to our format and filter
		for _, orgProject := range orgProjectsResp.Project {
			// Apply filters
			if projectID != 0 && orgProject.ProjectId != projectID {
				continue
			}
			if name != "" && orgProject.Name != name {
				continue
			}

			// Get project with tags from our repository
			projectWithTags, err := p.repository.ListProject(ctx, 0, orgProject.ProjectId, "")
			if err != nil {
				p.logger.Warnf(ctx, "Failed to get project details for project %d: %v", orgProject.ProjectId, err)
				continue
			}

			if projectWithTags != nil && len(*projectWithTags) > 0 {
				allProjects = append(allProjects, &(*projectWithTags)[0])
			}
		}
	}

	return allProjects, nil
}

// getUserOrganizations gets all organization IDs where the user has any role
func (p *ProjectService) getUserOrganizations(ctx context.Context, userID uint32) ([]uint32, error) {
	// Get all organizations (we'll filter by user role later)
	allOrgsResp, err := p.organizationClient.ListOrganization(ctx, &organization.ListOrganizationRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	var userOrgs []uint32

	// Check each organization to see if user has any role
	for _, org := range allOrgsResp.Organization {
		hasRole, err := p.userHasRoleInOrganization(ctx, userID, org.OrganizationId)
		if err != nil {
			p.logger.Warnf(ctx, "Failed to check user role in organization %d: %v", org.OrganizationId, err)
			continue
		}

		if hasRole {
			userOrgs = append(userOrgs, org.OrganizationId)
		}
	}

	return userOrgs, nil
}

// userHasRoleInOrganization checks if user has any role in the organization
func (p *ProjectService) userHasRoleInOrganization(ctx context.Context, userID, organizationID uint32) (bool, error) {
	// Try a simple organization action to see if user has any role
	authResp, err := p.organizationIamClient.IsAuthorizedOrganization(ctx, &organization_iam.IsAuthorizedOrganizationRequest{
		UserId:         userID,
		OrganizationId: organizationID,
		ActionName:     "organization/list", // Simple read action
	})
	if err != nil {
		// If there's an error, assume no role
		return false, nil
	}

	return authResp.Ok, nil
}

func convertProject(p *model.Project) *project.Project {
	return &project.Project{
		ProjectId: p.ProjectID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: p.UpdatedAt.Unix(),
	}
}

func (p *ProjectService) CreateProject(ctx context.Context, req *project.CreateProjectRequest) (*project.CreateProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	pr, err := p.repository.CreateProject(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if err := p.createDefaultRole(ctx, req.UserId, pr.ProjectID); err != nil {
		return nil, err
	}
	p.logger.Infof(ctx, "Project created: owner=%d, project=%+v", req.UserId, pr)

	return &project.CreateProjectResponse{Project: convertProject(pr)}, nil
}

func (p *ProjectService) UpdateProject(ctx context.Context, req *project.UpdateProjectRequest) (*project.UpdateProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	pr, err := p.repository.UpdateProject(ctx, req.ProjectId, req.Name)
	if err != nil {
		return nil, err
	}
	return &project.UpdateProjectResponse{Project: convertProject(pr)}, nil
}

func (p *ProjectService) DeleteProject(ctx context.Context, req *project.DeleteProjectRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := p.repository.DeleteProject(ctx, req.ProjectId); err != nil {
		return nil, err
	}
	p.logger.Infof(ctx, "Project deleted: project=%+v", req.ProjectId)

	return &empty.Empty{}, nil
}

func (p *ProjectService) IsActive(ctx context.Context, req *project.IsActiveRequest) (*project.IsActiveResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	active, err := p.isActiveProject(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return &project.IsActiveResponse{Active: active}, nil
}

func (p *ProjectService) createDefaultRole(ctx context.Context, ownerUserID, projectID uint32) error {
	projectAdmin := "project-admin"
	projectViewer := "project-viewer"
	findingEditor := "finding-editor"
	viewerActionPtn := "get|list|is-admin|put-alert-first-viewed-at"

	for name, actionPtn := range map[string]string{
		projectAdmin:  ".*",
		projectViewer: viewerActionPtn,
		findingEditor: viewerActionPtn + "|^finding/.+|^alert/.+",
	} {
		policy, err := p.iamClient.PutPolicy(ctx, &iam.PutPolicyRequest{
			ProjectId: projectID,
			Policy: &iam.PolicyForUpsert{
				Name:        name,
				ProjectId:   projectID,
				ActionPtn:   actionPtn,
				ResourcePtn: ".*",
			},
		})
		if err != nil {
			return fmt.Errorf("could not put %s-policy, err=%w", name, err)
		}
		role, err := p.iamClient.PutRole(ctx, &iam.PutRoleRequest{
			ProjectId: projectID,
			Role: &iam.RoleForUpsert{
				Name:      name + "-role",
				ProjectId: projectID,
			},
		})
		if err != nil {
			return fmt.Errorf("could not put %s-role, err=%w", name, err)
		}
		if _, err := p.iamClient.AttachPolicy(ctx, &iam.AttachPolicyRequest{
			ProjectId: projectID,
			RoleId:    role.Role.RoleId,
			PolicyId:  policy.Policy.PolicyId,
		}); err != nil {
			return fmt.Errorf("could not attach %s-policy to %s-role, err=%w", name, name, err)
		}
		if name == projectAdmin {
			if _, err := p.iamClient.AttachRole(ctx, &iam.AttachRoleRequest{
				ProjectId: projectID,
				UserId:    ownerUserID,
				RoleId:    role.Role.RoleId,
			}); err != nil {
				return fmt.Errorf("could not attach default %s-role to project owner, err=%w", name, err)
			}
		}
	}
	return nil
}

func (p *ProjectService) isActiveProject(ctx context.Context, projectID uint32) (bool, error) {
	projects, err := p.repository.ListProject(ctx, 0, projectID, "")
	if err != nil {
		return false, err
	}
	if len(*projects) == 0 {
		return false, nil
	}
	resp, err := p.iamClient.ListUser(ctx, &iam.ListUserRequest{
		ProjectId: projectID,
		Activated: true,
	})
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, nil
	}
	return len(resp.UserId) > 0, nil
}

func (p *ProjectService) CleanProject(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	if err := p.repository.CleanWithNoProject(ctx); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
