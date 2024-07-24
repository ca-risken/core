package iam

import (
	"context"
	"errors"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"gorm.io/gorm"
)

func (i *IAMService) ListUser(ctx context.Context, req *iam.ListUserRequest) (*iam.ListUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListUser(ctx, req.Activated, req.ProjectId, req.Name, req.UserId, req.Admin, req.UserIdpKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &iam.ListUserResponse{}, nil
		}
		return nil, err
	}
	ids := []uint32{}
	for _, u := range *list {
		ids = append(ids, u.UserID)
	}
	return &iam.ListUserResponse{UserId: ids}, nil
}

func (i *IAMService) GetUser(ctx context.Context, req *iam.GetUserRequest) (*iam.GetUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	user, err := i.repository.GetUser(ctx, req.UserId, req.Sub, req.UserIdpKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			i.logger.Infof(ctx, "[GetUser]User not found: GetUserRequest=%+v", req)
			return &iam.GetUserResponse{}, nil
		}
		return nil, err
	}
	return &iam.GetUserResponse{User: convertUser(user)}, nil
}

func (i *IAMService) PutUser(ctx context.Context, req *iam.PutUserRequest) (*iam.PutUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := i.repository.GetUserBySub(ctx, req.User.Sub)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, err
	}
	// check first user
	isFisrstUser := false
	if noRecord {
		users, err := i.repository.GetActiveUserCount(ctx)
		if err != nil {
			return nil, err
		}
		isFisrstUser = users == nil || *users == 0
	}
	i.logger.Debugf(ctx, "isFisrstUser: %t", isFisrstUser)

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var userID uint32
	if !noRecord {
		userID = savedData.UserID
	}
	u := &model.User{
		UserID:     userID,
		Sub:        req.User.Sub,
		Name:       req.User.Name,
		UserIdpKey: req.User.UserIdpKey,
		Activated:  req.User.Activated,
	}
	var registerdData *model.User
	// 登録済みユーザーの場合、update
	if userID != 0 {
		registerdData, err = i.repository.PutUser(ctx, u)
		if err != nil {
			return nil, err
		}
	} else {
		registerdData, err = i.repository.CreateUser(ctx, u)
		if err != nil {
			return nil, err
		}
	}

	if isFisrstUser {
		// attach admin roles
		if err := i.repository.AttachAllAdminRole(ctx, registerdData.UserID); err != nil {
			return nil, err
		}
		i.logger.Infof(ctx, "Attach admin role for first user, user_id=%d", registerdData.UserID)
	}
	// 新規ユーザーの場合、user_reservedからロールの追加
	if userID == 0 {
		if err := i.AttachRoleByUserReserved(ctx, registerdData.UserID, registerdData.UserIdpKey); err != nil {
			return nil, err
		}
	}

	return &iam.PutUserResponse{User: convertUser(registerdData)}, nil
}

func convertUser(u *model.User) *iam.User {
	return &iam.User{
		UserId:     u.UserID,
		Sub:        u.Sub,
		Name:       u.Name,
		UserIdpKey: u.UserIdpKey,
		Activated:  u.Activated,
		CreatedAt:  u.CreatedAt.Unix(),
		UpdatedAt:  u.UpdatedAt.Unix(),
	}
}
