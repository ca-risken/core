package main

import (
	"context"
	"errors"

	"github.com/ca-risken/core/proto/project"
	"github.com/ca-risken/core/src/project/model"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func convertProjectWithTag(p *projectWithTag) *project.Project {
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

func (p *projectService) ListProject(ctx context.Context, req *project.ListProjectRequest) (*project.ListProjectResponse, error) {
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

func (p *projectService) CreateProject(ctx context.Context, req *project.CreateProjectRequest) (*project.CreateProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	pr, err := p.repository.CreateProject(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if err := p.iamClient.CreateDefaultRole(ctx, req.UserId, pr.ProjectID); err != nil {
		return nil, err
	}
	appLogger.Infof("Project created: owner=%d, project=%+v", req.UserId, pr)
	return &project.CreateProjectResponse{Project: convertProject(pr)}, nil
}

func (p *projectService) UpdateProject(ctx context.Context, req *project.UpdateProjectRequest) (*project.UpdateProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	pr, err := p.repository.UpdateProject(ctx, req.ProjectId, req.Name)
	if err != nil {
		return nil, err
	}
	return &project.UpdateProjectResponse{Project: convertProject(pr)}, nil
}

func (p *projectService) DeleteProject(ctx context.Context, req *project.DeleteProjectRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := p.iamClient.DeleteAllProjectRole(ctx, req.ProjectId); err != nil {
		return nil, err
	}
	appLogger.Infof("Project deleted: project=%+v", req.ProjectId)
	return &empty.Empty{}, nil
}
