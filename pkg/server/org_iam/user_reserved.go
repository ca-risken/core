package org_iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/org_iam"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (i *OrgIAMService) ListOrgUserReserved(ctx context.Context, req *org_iam.ListOrgUserReservedRequest) (*org_iam.ListOrgUserReservedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListOrgUserReserved(ctx, req.OrganizationId, req.UserIdpKey)
	if err != nil {
		return nil, err
	}
	var userReserved []*org_iam.OrgUserReserved
	for _, l := range list {
		userReserved = append(userReserved, convertOrgUserReserved(l))
	}
	return &org_iam.ListOrgUserReservedResponse{UserReserved: userReserved}, nil
}

func (i *OrgIAMService) PutOrgUserReserved(ctx context.Context, req *org_iam.PutOrgUserReservedRequest) (*org_iam.PutOrgUserReservedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	user, err := i.iamClient.GetUserByUserIdpKey(ctx, &iam.GetUserByUserIdpKeyRequest{UserIdpKey: req.UserIdpKey})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if user.User != nil && user.User.UserId != 0 {
		return nil, errors.New("user already exists")
	}
	_, err = i.repository.GetOrgRole(ctx, req.OrganizationId, req.RoleId)
	if err != nil {
		return nil, fmt.Errorf("role does not exist in the organization, roleID: %v, organizationID: %v", req.RoleId, req.OrganizationId)
	}
	putData := &model.OrganizationUserReserved{
		ReservedID: req.ReservedId,
		UserIdpKey: req.UserIdpKey,
		RoleID:     req.RoleId,
	}
	res, err := i.repository.PutOrgUserReserved(ctx, putData)
	if err != nil {
		return nil, err
	}
	return &org_iam.PutOrgUserReservedResponse{UserReserved: convertOrgUserReserved(res)}, nil
}

func (i *OrgIAMService) DeleteOrgUserReserved(ctx context.Context, req *org_iam.DeleteOrgUserReservedRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DeleteOrgUserReserved(ctx, req.OrganizationId, req.ReservedId); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func convertOrgUserReserved(u *model.OrganizationUserReserved) *org_iam.OrgUserReserved {
	return &org_iam.OrgUserReserved{
		ReservedId: u.ReservedID,
		UserIdpKey: u.UserIdpKey,
		RoleId:     u.RoleID,
		CreatedAt:  u.CreatedAt.Unix(),
		UpdatedAt:  u.UpdatedAt.Unix(),
	}
}
