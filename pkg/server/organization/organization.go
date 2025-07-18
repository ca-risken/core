package organization

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/organization"
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/ca-risken/core/proto/project"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (o *OrganizationService) ListOrganization(ctx context.Context, req *organization.ListOrganizationRequest) (*organization.ListOrganizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := o.repository.ListOrganization(ctx, req.OrganizationId, req.Name, req.UserId, req.ProjectId)
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
	if err := o.createDefaultRole(ctx, req.UserId, org.OrganizationID); err != nil {
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
	if req.OrganizationId == 0 && req.ProjectId == 0 {
		return nil, errors.New("at least one of organizationID or projectID must be specified")
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
	exists, err := o.repository.ExistsOrganizationProject(ctx, req.OrganizationId, req.ProjectId)
	if err != nil {
		return nil, err
	}
	if exists && req.Status != organization.OrganizationInvitationStatus_ACCEPTED {
		return nil, errors.New("organization is already associated with the project")
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
	if err := o.removeOrganizationProject(ctx, req.OrganizationId, req.ProjectId); err != nil {
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
	if invitation.Status == organization.OrganizationInvitationStatus_REJECTED.String() {
		if err := o.removeOrganizationProject(ctx, req.OrganizationId, req.ProjectId); err != nil {
			return nil, err
		}
	}
	return &organization.ReplyOrganizationInvitationResponse{}, nil
}

func (o *OrganizationService) removeOrganizationProject(ctx context.Context, organizationID, projectID uint32) error {
	exists, err := o.repository.ExistsOrganizationProject(ctx, organizationID, projectID)
	if err != nil {
		return err
	}
	if exists {
		if err := o.repository.RemoveProjectsInOrganization(ctx, organizationID, projectID); err != nil {
			return err
		}
		o.logger.Infof(ctx, "OrganizationProject removed due to invitation deletion: organization_id=%d, project_id=%d", organizationID, projectID)
	}
	return nil
}

func (o *OrganizationService) createDefaultRole(ctx context.Context, ownerUserID, organizationID uint32) error {
	organizationAdmin := "organization-admin"
	organizationViewer := "organization-viewer"
	viewerActionPtn := "get|list|is-admin|put-alert-first-viewed-at"

	for name, actionPtn := range map[string]string{
		organizationAdmin:  ".*",
		organizationViewer: viewerActionPtn,
	} {
		policy, err := o.organizationIamClient.PutOrganizationPolicy(ctx, &organization_iam.PutOrganizationPolicyRequest{
			OrganizationId: organizationID,
			Name:           name,
			ActionPtn:      actionPtn,
		})
		if err != nil {
			return fmt.Errorf("could not put %s-policy, err=%w", name, err)
		}
		role, err := o.organizationIamClient.PutOrganizationRole(ctx, &organization_iam.PutOrganizationRoleRequest{
			OrganizationId: organizationID,
			Name:           name + "-role",
		})
		if err != nil {
			return fmt.Errorf("could not put %s-role, err=%w", name, err)
		}
		if _, err := o.organizationIamClient.AttachOrganizationPolicy(ctx, &organization_iam.AttachOrganizationPolicyRequest{
			OrganizationId: organizationID,
			RoleId:         role.Role.RoleId,
			PolicyId:       policy.Policy.PolicyId,
		}); err != nil {
			return fmt.Errorf("could not attach %s-policy to %s-role, err=%w", name, name, err)
		}
		if name == organizationAdmin {
			if _, err := o.organizationIamClient.AttachOrganizationRole(ctx, &organization_iam.AttachOrganizationRoleRequest{
				OrganizationId: organizationID,
				UserId:         ownerUserID,
				RoleId:         role.Role.RoleId,
			}); err != nil {
				return fmt.Errorf("could not attach default %s-role to organization owner, err=%w", name, err)
			}
		}
	}
	return nil
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
