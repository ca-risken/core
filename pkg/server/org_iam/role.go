package org_iam

import (
	"context"
	"errors"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/org_iam"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (i *OrgIAMService) ListOrgRole(ctx context.Context, req *org_iam.ListOrgRoleRequest) (*org_iam.ListOrgRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListOrgRole(ctx, req.OrganizationId, req.Name, req.UserId, req.AccessTokenId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_iam.ListOrgRoleResponse{}, nil
		}
		return nil, err
	}
	ids := []uint32{}
	for _, r := range list {
		ids = append(ids, r.RoleID)
	}
	return &org_iam.ListOrgRoleResponse{RoleId: ids}, nil
}

func (i *OrgIAMService) GetOrgRole(ctx context.Context, req *org_iam.GetOrgRoleRequest) (*org_iam.GetOrgRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	role, err := i.repository.GetOrgRole(ctx, req.OrganizationId, req.RoleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_iam.GetOrgRoleResponse{}, nil
		}
		return nil, err
	}
	return &org_iam.GetOrgRoleResponse{Role: convertOrgRole(role)}, nil
}

func convertOrgRole(r *model.OrganizationRole) *org_iam.OrgRole {
	return &org_iam.OrgRole{
		RoleId:         r.RoleID,
		Name:           r.Name,
		OrganizationId: r.OrganizationID,
		CreatedAt:      r.CreatedAt.Unix(),
		UpdatedAt:      r.UpdatedAt.Unix(),
	}
}

func (i *OrgIAMService) PutOrgRole(ctx context.Context, req *org_iam.PutOrgRoleRequest) (*org_iam.PutOrgRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := i.repository.GetOrgRoleByName(ctx, req.OrganizationId, req.Name)
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
	registerdData, err := i.repository.PutOrgRole(ctx, r)
	if err != nil {
		return nil, err
	}
	return &org_iam.PutOrgRoleResponse{Role: convertOrgRole(registerdData)}, nil
}

func (i *OrgIAMService) DeleteOrgRole(ctx context.Context, req *org_iam.DeleteOrgRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DeleteOrgRole(ctx, req.OrganizationId, req.RoleId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *OrgIAMService) AttachOrgRole(ctx context.Context, req *org_iam.AttachOrgRoleRequest) (*org_iam.AttachOrgRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	role, err := s.repository.AttachOrgRole(ctx, req.OrganizationId, req.RoleId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &org_iam.AttachOrgRoleResponse{Role: convertOrgRole(role)}, nil
}

func (s *OrgIAMService) DetachOrgRole(ctx context.Context, req *org_iam.DetachOrgRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DetachOrgRole(ctx, req.OrganizationId, req.RoleId, req.UserId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
func (i *OrgIAMService) AttachOrgRoleByOrgUserReserved(ctx context.Context, req *org_iam.AttachOrgRoleByOrgUserReservedRequest) (*org_iam.AttachOrgRoleByOrgUserReservedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	userReserved, err := i.repository.ListOrgUserReservedWithOrganizationID(ctx, req.UserIdpKey)
	if err != nil {
		return nil, err
	}
	for _, u := range *userReserved {
		_, err := i.repository.AttachOrgRole(ctx, u.OrganizationID, u.RoleID, req.UserId)
		if err != nil {
			return nil, err
		}
		if err := i.repository.DeleteOrgUserReserved(ctx, u.OrganizationID, u.ReservedID); err != nil {
			return nil, err
		}
	}
	return &org_iam.AttachOrgRoleByOrgUserReservedResponse{}, nil
}
