package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/golang/protobuf/ptypes/empty"
)

func (i *iamService) ListRole(ctx context.Context, req *iam.ListRoleRequest) (*iam.ListRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *iamService) GetRole(ctx context.Context, req *iam.GetRoleRequest) (*iam.GetRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *iamService) PutRole(ctx context.Context, req *iam.PutRoleRequest) (*iam.PutRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *iamService) DeleteRole(ctx context.Context, req *iam.DeleteRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *iamService) AttachRole(ctx context.Context, req *iam.AttachRoleRequest) (*iam.AttachRoleResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *iamService) DetachRole(ctx context.Context, req *iam.DetachRoleRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}
