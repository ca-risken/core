package organization

import (
	"context"
	"errors"
	"strings"

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

func (o *OrganizationService) PutOrganizationInvitation(ctx context.Context, req *organization.PutOrganizationInvitationRequest) (*organization.PutOrganizationInvitationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	invitation, err := o.repository.PutOrganizationInvitation(ctx, req.OrganizationId, req.ProjectId, req.Status.String())
	if err != nil {
		return nil, err
	}
	return &organization.PutOrganizationInvitationResponse{OrganizationInvitation: convertOrganizationInvitation(invitation)}, nil
}

func (o *OrganizationService) DeleteOrganizationInvitation(ctx context.Context, req *organization.DeleteOrganizationInvitationRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := o.repository.DeleteOrganizationInvitation(ctx, req.OrganizationId, req.ProjectId); err != nil {
		return nil, err
	}
	o.logger.Infof(ctx, "Organization invitation deleted: organization_id=%d, project_id=%d", req.OrganizationId, req.ProjectId)
	return &empty.Empty{}, nil
}

func (o *OrganizationService) ReplyOrganizationInvitation(ctx context.Context, req *organization.ReplyOrganizationInvitationRequest) (*organization.ReplyOrganizationInvitationResponse, error) {
	var orgProject *model.OrganizationProject
	if err := req.Validate(); err != nil {
		return nil, err
	}
	invitation, err := o.repository.PutOrganizationInvitation(ctx, req.OrganizationId, req.ProjectId, req.Status.String())
	if err != nil {
		return nil, err
	}
	if invitation.Status == organization.OrganizationInvitationStatus_ACCEPTED.String() {
		orgProject, err = o.repository.PutOrganizationProject(ctx, req.OrganizationId, req.ProjectId)
		if err != nil {
			return nil, err
		}
		return &organization.ReplyOrganizationInvitationResponse{OrganizationProject: convertOrganizationProject(orgProject)}, nil
	}
	return &organization.ReplyOrganizationInvitationResponse{}, nil
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
		Status:         getOrganizationInvitationStatus(oi.Status),
		CreatedAt:      oi.CreatedAt.Unix(),
		UpdatedAt:      oi.UpdatedAt.Unix(),
	}
}

func convertOrganizationProject(op *model.OrganizationProject) *organization.OrganizationProject {
	return &organization.OrganizationProject{
		OrganizationId: op.OrganizationID,
		ProjectId:      op.ProjectID,
		CreatedAt:      op.CreatedAt.Unix(),
		UpdatedAt:      op.UpdatedAt.Unix(),
	}
}

func getOrganizationInvitationStatus(s string) organization.OrganizationInvitationStatus {
	statusKey := strings.ToUpper(s)
	if _, ok := organization.OrganizationInvitationStatus_value[statusKey]; !ok {
		return organization.OrganizationInvitationStatus_UNKNOWN
	}
	switch statusKey {
	case organization.OrganizationInvitationStatus_PENDING.String():
		return organization.OrganizationInvitationStatus_PENDING
	case organization.OrganizationInvitationStatus_ACCEPTED.String():
		return organization.OrganizationInvitationStatus_ACCEPTED
	case organization.OrganizationInvitationStatus_REJECTED.String():
		return organization.OrganizationInvitationStatus_REJECTED
	default:
		return organization.OrganizationInvitationStatus_UNKNOWN
	}
}
