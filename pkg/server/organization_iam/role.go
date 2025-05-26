package organization_iam

import (
	"context"
	"errors"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (i *OrganizationIAMService) ListOrganizationRole(ctx context.Context, req *organization_iam.ListOrganizationRoleRequest) (*organization_iam.ListOrganizationRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListOrganizationRole(ctx, req.OrganizationId, req.Name, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.ListOrganizationRoleResponse{}, nil
		}
		return nil, err
	}
	ids := []uint32{}
	for _, r := range list {
		ids = append(ids, r.RoleID)
	}
	return &organization_iam.ListOrganizationRoleResponse{RoleId: ids}, nil
}

func (i *OrganizationIAMService) GetOrganizationRole(ctx context.Context, req *organization_iam.GetOrganizationRoleRequest) (*organization_iam.GetOrganizationRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	role, err := i.repository.GetOrganizationRole(ctx, req.OrganizationId, req.RoleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.GetOrganizationRoleResponse{}, nil
		}
		return nil, err
	}
	return &organization_iam.GetOrganizationRoleResponse{Role: convertOrganizationRole(role)}, nil
}

func convertOrganizationRole(r *model.OrganizationRole) *organization_iam.OrganizationRole {
	return &organization_iam.OrganizationRole{
		RoleId:         r.RoleID,
		Name:           r.Name,
		OrganizationId: r.OrganizationID,
		CreatedAt:      r.CreatedAt.Unix(),
		UpdatedAt:      r.UpdatedAt.Unix(),
	}
}

func (i *OrganizationIAMService) PutOrganizationRole(ctx context.Context, req *organization_iam.PutOrganizationRoleRequest) (*organization_iam.PutOrganizationRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := i.repository.GetOrganizationRoleByName(ctx, req.OrganizationId, req.Name)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, err
	}
	var roleID uint32
	if !noRecord {
		roleID = savedData.RoleID
	}
	r := &model.OrganizationRole{
		RoleID:         roleID,
		Name:           req.Name,
		OrganizationID: req.OrganizationId,
	}
	registerdData, err := i.repository.PutOrganizationRole(ctx, r)
	if err != nil {
		return nil, err
	}
	return &organization_iam.PutOrganizationRoleResponse{Role: convertOrganizationRole(registerdData)}, nil
}

func (i *OrganizationIAMService) DeleteOrganizationRole(ctx context.Context, req *organization_iam.DeleteOrganizationRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DeleteOrganizationRole(ctx, req.OrganizationId, req.RoleId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
