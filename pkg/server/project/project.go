package project

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
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
				// CreatedAt: t.CreatedAt.Unix(), // Reduce the API response size
				// UpdatedAt: t.UpdatedAt.Unix(),
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
	list, err := p.repository.ListProject(ctx, req.UserId, req.ProjectId, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &project.ListProjectResponse{}, nil
		}
		return nil, err
	}
	var prs []*project.Project
	for _, pr := range *list {
		prs = append(prs, convertProjectWithTag(&pr))
	}
	return &project.ListProjectResponse{Project: prs}, nil
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
	appLogger.Infof("Project created: owner=%d, project=%+v", req.UserId, pr)
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
	if err := p.deleteAllProjectRole(ctx, req.ProjectId); err != nil {
		return nil, err
	}
	appLogger.Infof("Project deleted: project=%+v", req.ProjectId)
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
	policy, err := p.iamClient.PutPolicy(ctx, &iam.PutPolicyRequest{
		ProjectId: projectID,
		Policy: &iam.PolicyForUpsert{
			Name:        "project-admin",
			ProjectId:   projectID,
			ActionPtn:   ".*",
			ResourcePtn: ".*",
		},
	})
	if err != nil {
		return fmt.Errorf("Could not put default policy, err=%+v", err)
	}
	role, err := p.iamClient.PutRole(ctx, &iam.PutRoleRequest{
		ProjectId: projectID,
		Role: &iam.RoleForUpsert{
			Name:      "project-admin-role",
			ProjectId: projectID,
		},
	})
	if err != nil {
		return fmt.Errorf("Could not put project-admin-role, err=%+v", err)
	}
	if _, err := p.iamClient.AttachPolicy(ctx, &iam.AttachPolicyRequest{
		ProjectId: projectID,
		RoleId:    role.Role.RoleId,
		PolicyId:  policy.Policy.PolicyId,
	}); err != nil {
		return fmt.Errorf("Could not attach default policy, err=%+v", err)
	}
	if _, err := p.iamClient.AttachRole(ctx, &iam.AttachRoleRequest{
		ProjectId: projectID,
		UserId:    ownerUserID,
		RoleId:    role.Role.RoleId,
	}); err != nil {
		return fmt.Errorf("Could not attach default role, err=%+v", err)
	}
	return nil
}

func (p *ProjectService) deleteAllProjectRole(ctx context.Context, projectID uint32) error {
	list, err := p.iamClient.ListRole(ctx, &iam.ListRoleRequest{ProjectId: projectID})
	if err != nil {
		return err
	}
	for _, roleID := range list.RoleId {
		if _, err := p.iamClient.DeleteRole(ctx, &iam.DeleteRoleRequest{
			ProjectId: projectID,
			RoleId:    roleID,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (p *ProjectService) isActiveProject(ctx context.Context, projectID uint32) (bool, error) {
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