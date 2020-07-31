package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/project"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type projectService struct {
	repository projectRepository
	iamClient  iamService
}

type projectConfig struct {
	Port  string `default:"8003"`
	Debug bool   `default:"false"`
}

func newProjectService() project.ProjectServiceServer {
	var conf projectConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		appLogger.Fatalf("project config load error: err=%+v", err)
	}
	if conf.Debug {
		appLogger.SetLevel(logrus.DebugLevel)
	}
	return &projectService{
		repository: newProjectRepository(),
		iamClient:  newIAMService(),
	}
}

func (p *projectService) ListProject(ctx context.Context, req *project.ListProjectRequest) (*project.ListProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := p.repository.ListProject(req.UserId, req.ProjectId, req.Name)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &project.ListProjectResponse{}, nil
		}
		return nil, err
	}
	var prs []*project.Project
	for _, pr := range *list {
		prs = append(prs, convertProject(&pr))
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
	pr, err := p.repository.CreateProject(req.Name)
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
	pr, err := p.repository.UpdateProject(req.ProjectId, req.Name)
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
	return &empty.Empty{}, nil
}