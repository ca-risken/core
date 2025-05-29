package organization

import (
	"context"
	"errors"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/project"
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
	for _, org := range list {
		orgs = append(orgs, convertOrganization(org))
	}
	return &organization.ListOrganizationResponse{Organization: orgs}, nil
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

func (o *OrganizationService) ListProjectsInOrganization(ctx context.Context, req *organization.ListProjectsInOrganizationRequest) (*organization.ListProjectsInOrganizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	projects, err := o.repository.ListProjectsInOrganization(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	var result []*project.Project
	for _, p := range projects {
		result = append(result, &project.Project{
			ProjectId: p.ProjectID,
			Name:      p.Name,
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		})
	}
	return &organization.ListProjectsInOrganizationResponse{Project: result}, nil
}

func (o *OrganizationService) RemoveProjectsInOrganization(ctx context.Context, req *organization.RemoveProjectsInOrganizationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := o.repository.RemoveProjectsInOrganization(ctx, req.OrganizationId, req.ProjectId); err != nil {
		return nil, err
	}
	o.logger.Infof(ctx, "Projects removed from organization: organization_id=%d, project_id=%d", req.OrganizationId, req.ProjectId)
	return &empty.Empty{}, nil
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
