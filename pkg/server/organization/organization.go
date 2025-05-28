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

func (o *OrganizationService) ListProjectsByOrganization(ctx context.Context, req *organization.ListProjectsByOrganizationRequest) (*organization.ListProjectsByOrganizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	projects, err := o.repository.ListProjectsByOrganization(ctx, req.OrganizationId)
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
	return &organization.ListProjectsByOrganizationResponse{Project: result}, nil
}

func (o *OrganizationService) AddProjects(ctx context.Context, req *organization.AddProjectsRequest) (*organization.AddProjectsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	orgProject, err := o.repository.AddProjects(ctx, req.OrganizationId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	return &organization.AddProjectsResponse{OrganizationProject: &organization.OrganizationProject{
		OrganizationId: orgProject.OrganizationID,
		ProjectId:      orgProject.ProjectID,
		CreatedAt:      orgProject.CreatedAt.Unix(),
		UpdatedAt:      orgProject.UpdatedAt.Unix(),
	}}, nil
}

func (o *OrganizationService) ListOrganizationInvitation(ctx context.Context, req *organization.ListOrganizationInvitationRequest) (*organization.ListOrganizationInvitationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	invitations, err := o.repository.ListOrganizationInvitation(ctx, req.OrganizationId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	var result []*organization.OrganizationInvitation
	for _, invitation := range invitations {
		result = append(result, convertOrganizationInvitation(invitation))
	}
	return &organization.ListOrganizationInvitationResponse{OrganizationInvitations: result}, nil
}

func (o *OrganizationService) CreateOrganizationInvitation(ctx context.Context, req *organization.CreateOrganizationInvitationRequest) (*organization.CreateOrganizationInvitationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	invitation, err := o.repository.CreateOrganizationInvitation(ctx, req.OrganizationId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	o.logger.Infof(ctx, "Organization invitation created: organization_id=%d, project_id=%d", req.OrganizationId, req.ProjectId)
	return &organization.CreateOrganizationInvitationResponse{OrganizationInvitation: convertOrganizationInvitation(invitation)}, nil
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

func convertOrganizationInvitation(oi *model.OrganizationInvitation) *organization.OrganizationInvitation {
	return &organization.OrganizationInvitation{
		OrganizationId: oi.OrganizationID,
		ProjectId:      oi.ProjectID,
		Status:         oi.Status,
		CreatedAt:      oi.CreatedAt.Unix(),
		UpdatedAt:      oi.UpdatedAt.Unix(),
	}
}
