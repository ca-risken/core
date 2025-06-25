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
