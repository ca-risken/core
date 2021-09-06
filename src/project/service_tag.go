package main

import (
	"context"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/project"
	"github.com/golang/protobuf/ptypes/empty"
)

func convertProjectTag(p *model.ProjectTag) *project.ProjectTag {
	return &project.ProjectTag{
		ProjectId: p.ProjectID,
		Tag:       p.Tag,
		Color:     p.Color,
		CreatedAt: p.CreatedAt.Unix(),
		UpdatedAt: p.UpdatedAt.Unix(),
	}
}

func (p *projectService) TagProject(ctx context.Context, req *project.TagProjectRequest) (*project.TagProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	tag, err := p.repository.TagProject(ctx, req.ProjectId, req.Tag, req.Color)
	if err != nil {
		return nil, err
	}
	return &project.TagProjectResponse{ProjectTag: convertProjectTag(tag)}, nil
}

func (p *projectService) UntagProject(ctx context.Context, req *project.UntagProjectRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := p.repository.UntagProject(ctx, req.ProjectId, req.Tag); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
