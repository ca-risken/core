package organization_iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization_iam"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (i *OrganizationIAMService) ListOrganizationUserReserved(ctx context.Context, req *organization_iam.ListOrganizationUserReservedRequest) (*organization_iam.ListOrganizationUserReservedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListOrganizationUserReserved(ctx, req.OrganizationId, req.UserIdpKey)
	if err != nil {
		return nil, err
	}
	var userReserved []*organization_iam.OrganizationUserReserved
	for _, l := range list {
		userReserved = append(userReserved, convertOrganizationUserReserved(l))
	}
	return &organization_iam.ListOrganizationUserReservedResponse{UserReserved: userReserved}, nil
}

func (i *OrganizationIAMService) PutOrganizationUserReserved(ctx context.Context, req *organization_iam.PutOrganizationUserReservedRequest) (*organization_iam.PutOrganizationUserReservedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	user, err := i.iamClient.GetUser(ctx, &iam.GetUserRequest{UserIdpKey: req.UserIdpKey})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if user.User != nil && user.User.UserId != 0 {
		return nil, errors.New("user already exists")
	}
	_, err = i.repository.GetOrganizationRole(ctx, req.OrganizationId, req.RoleId)
	if err != nil {
		return nil, fmt.Errorf("role dose not exist in the orgazniation, roleID: %v, organizationID: %v", req.RoleId, req.OrganizationId)
	}
	putData := &model.OrganizationUserReserved{
		ReservedID: req.ReservedId,
		UserIdpKey: req.UserIdpKey,
		RoleID:     req.RoleId,
	}
	res, err := i.repository.PutOrganizationUserReserved(ctx, putData)
	if err != nil {
		return nil, err
	}
	return &organization_iam.PutOrganizationUserReservedResponse{UserReserved: convertOrganizationUserReserved(res)}, nil
}

func (i *OrganizationIAMService) DeleteOrganizationUserReserved(ctx context.Context, req *organization_iam.DeleteOrganizationUserReservedRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DeleteOrganizationUserReserved(ctx, req.OrganizationId, req.ReservedId); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func convertOrganizationUserReserved(u *model.OrganizationUserReserved) *organization_iam.OrganizationUserReserved {
	return &organization_iam.OrganizationUserReserved{
		ReservedId: u.ReservedID,
		UserIdpKey: u.UserIdpKey,
		RoleId:     u.RoleID,
		CreatedAt:  u.CreatedAt.Unix(),
		UpdatedAt:  u.UpdatedAt.Unix(),
	}
}
