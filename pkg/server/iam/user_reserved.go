package iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (i *IAMService) ListUserReserved(ctx context.Context, req *iam.ListUserReservedRequest) (*iam.ListUserReservedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListUserReserved(ctx, req.ProjectId, req.UserIdpKey)
	if err != nil {
		return nil, err
	}
	var userReserved []*iam.UserReserved
	for _, l := range *list {
		userReserved = append(userReserved, convertUserReserved(&l))
	}
	return &iam.ListUserReservedResponse{UserReserved: userReserved}, nil
}

func (i *IAMService) PutUserReserved(ctx context.Context, req *iam.PutUserReservedRequest) (*iam.PutUserReservedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	user, err := i.repository.GetUserByUserIdpKey(ctx, req.UserReserved.UserIdpKey)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, err
	}
	// 既にユーザーが存在する場合、プロジェクトにユーザーを招待すれば良いので処理終了
	if user != nil && user.UserID != 0 {
		return nil, errors.New("already existing user. Invite user to your project")
	}

	// 割り当てる権限が対象のプロジェクトに存在するかを確認する
	// 存在しない場合、エラーを返す
	_, err = i.repository.GetRole(ctx, req.ProjectId, req.UserReserved.RoleId)
	if err != nil {
		return nil, fmt.Errorf("role does not exist in the project, roleID: %v, projectID: %v", req.UserReserved.RoleId, req.ProjectId)
	}
	putData := &model.UserReserved{
		ReservedID: req.UserReserved.ReservedId,
		UserIdpKey: req.UserReserved.UserIdpKey,
		RoleID:     req.UserReserved.RoleId,
	}
	registerdData, err := i.repository.PutUserReserved(ctx, putData)
	if err != nil {
		return nil, err
	}

	return &iam.PutUserReservedResponse{UserReserved: convertUserReserved(registerdData)}, nil
}

func (i *IAMService) DeleteUserReserved(ctx context.Context, req *iam.DeleteUserReservedRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	err := i.repository.DeleteUserReserved(ctx, req.ProjectId, req.ReservedId)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (i *IAMService) AttachRoleByUserReserved(ctx context.Context, userID uint32, userIdpKey string) error {
	userReserved, err := i.repository.ListUserReservedWithProjectID(ctx, userIdpKey)
	if err != nil {
		return err
	}
	for _, u := range *userReserved {
		_, err := i.repository.AttachRole(ctx, u.ProjectID, u.RoleID, userID)
		if err != nil {
			return err
		}
		if err := i.repository.DeleteUserReserved(ctx, u.ProjectID, u.ReservedID); err != nil {
			return err
		}
	}

	return nil
}

func convertUserReserved(u *model.UserReserved) *iam.UserReserved {
	return &iam.UserReserved{
		ReservedId: u.ReservedID,
		UserIdpKey: u.UserIdpKey,
		RoleId:     u.RoleID,
		CreatedAt:  u.CreatedAt.Unix(),
		UpdatedAt:  u.UpdatedAt.Unix(),
	}
}
