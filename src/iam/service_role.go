package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

func (i *iamService) ListRole(ctx context.Context, req *iam.ListRoleRequest) (*iam.ListRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListRole(req.ProjectId, req.Name)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &iam.ListRoleResponse{}, nil
		}
		return nil, err
	}
	ids := []uint32{}
	for _, r := range *list {
		ids = append(ids, r.RoleID)
	}
	return &iam.ListRoleResponse{RoleId: ids}, nil
}

func (i *iamService) GetRole(ctx context.Context, req *iam.GetRoleRequest) (*iam.GetRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	role, err := i.repository.GetRole(req.ProjectId, req.RoleId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &iam.GetRoleResponse{}, nil
		}
		return nil, err
	}
	return &iam.GetRoleResponse{Role: convertRole(role)}, nil
}

func convertRole(r *model.Role) *iam.Role {
	return &iam.Role{
		RoleId:    r.RoleID,
		Name:      r.Name,
		ProjectId: r.ProjectID,
		CreatedAt: r.CreatedAt.Unix(),
		UpdatedAt: r.UpdatedAt.Unix(),
	}
}

func (i *iamService) PutRole(ctx context.Context, req *iam.PutRoleRequest) (*iam.PutRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := i.repository.GetRoleByName(req.Role.ProjectId, req.Role.Name)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var roleID uint32
	if !noRecord {
		roleID = savedData.RoleID
	}
	r := &model.Role{
		RoleID:    roleID,
		Name:      req.Role.Name,
		ProjectID: req.Role.ProjectId,
	}

	// upsert
	registerdData, err := i.repository.PutRole(r)
	if err != nil {
		return nil, err
	}
	return &iam.PutRoleResponse{Role: convertRole(registerdData)}, nil
}

func (i *iamService) DeleteRole(ctx context.Context, req *iam.DeleteRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DeleteRole(req.ProjectId, req.RoleId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (i *iamService) AttachRole(ctx context.Context, req *iam.AttachRoleRequest) (*iam.AttachRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	ur, err := i.repository.AttachRole(req.ProjectId, req.RoleId, req.UserId)
	if err != nil {
		return nil, err
	}
	return &iam.AttachRoleResponse{UserRole: convertUserRole(ur)}, nil
}

func convertUserRole(ur *model.UserRole) *iam.UserRole {
	return &iam.UserRole{
		UserId:    ur.UserID,
		RoleId:    ur.RoleID,
		ProjectId: ur.ProjectID,
		CreatedAt: ur.CreatedAt.Unix(),
		UpdatedAt: ur.UpdatedAt.Unix(),
	}
}

func (i *iamService) DetachRole(ctx context.Context, req *iam.DetachRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DetachRole(req.ProjectId, req.RoleId, req.UserId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}