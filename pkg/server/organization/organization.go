package organization

import (
	"context"
	"errors"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/organization"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (o *OrganizationService) ListOrganization(ctx context.Context, req *organization.ListOrganizationRequest) (*organization.ListOrganizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := o.repository.ListOrganization(ctx, req.OrganizationId, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization.ListOrganizationResponse{}, nil
		}
		return nil, err
	}
	var orgs []*organization.Organization
	for _, org := range *list {
		orgs = append(orgs, convertOrganization(&org))
	}
	return &organization.ListOrganizationResponse{Organization: orgs}, nil
}

func convertOrganization(o *model.Organization) *organization.Organization {
	return &organization.Organization{
		OrganizationId: o.OrganizationID,
		Name:           o.Name,
		Description:    o.Description,
		CreatedAt:      o.CreatedAt.Unix(),
		UpdatedAt:      o.UpdatedAt.Unix(),
	}
}

func (o *OrganizationService) CreateOrganization(ctx context.Context, req *organization.CreateOrganizationRequest) (*organization.CreateOrganizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	org, err := o.repository.CreateOrganization(ctx, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	o.logger.Infof(ctx, "Organization created: organization=%+v", org)
	return &organization.CreateOrganizationResponse{Organization: convertOrganization(org)}, nil
}

func (o *OrganizationService) UpdateOrganization(ctx context.Context, req *organization.UpdateOrganizationRequest) (*organization.UpdateOrganizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	org, err := o.repository.UpdateOrganization(ctx, req.OrganizationId, req.Name, req.Description)
	if err != nil {
		return nil, err
	}
	return &organization.UpdateOrganizationResponse{Organization: convertOrganization(org)}, nil
}

func (o *OrganizationService) DeleteOrganization(ctx context.Context, req *organization.DeleteOrganizationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := o.repository.DeleteOrganization(ctx, req.OrganizationId); err != nil {
		return nil, err
	}
	o.logger.Infof(ctx, "Organization deleted: organization=%+v", req.OrganizationId)
	return &empty.Empty{}, nil
}

func convertOrganizationProject(p *model.OrganizationProject) *organization.OrganizationProject {
	return &organization.OrganizationProject{
		OrganizationId: p.OrganizationID,
		ProjectId:      p.ProjectID,
		CreatedAt:      p.CreatedAt.Unix(),
		UpdatedAt:      p.UpdatedAt.Unix(),
	}
}

func (o *OrganizationService) ListOrganizationProject(ctx context.Context, req *organization.ListOrganizationProjectRequest) (*organization.ListOrganizationProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := o.repository.ListOrganizationProject(ctx, req.OrganizationId, req.ProjectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization.ListOrganizationProjectResponse{}, nil
		}
		return nil, err
	}
	var projects []*organization.OrganizationProject
	for _, p := range *list {
		projects = append(projects, convertOrganizationProject(&p))
	}
	return &organization.ListOrganizationProjectResponse{OrganizationProject: projects}, nil
}

func (o *OrganizationService) GetOrganizationProject(ctx context.Context, req *organization.GetOrganizationProjectRequest) (*organization.GetOrganizationProjectResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	project, err := o.repository.GetOrganizationProject(ctx, req.OrganizationId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return &organization.GetOrganizationProjectResponse{OrganizationProject: convertOrganizationProject(project)}, nil
}
