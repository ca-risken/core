package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/jinzhu/gorm"
)

func (i *iamService) ListUser(ctx context.Context, req *iam.ListUserRequest) (*iam.ListUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListUser(req.Activated, req.ProjectId, req.Name, req.UserId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
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

func (i *iamService) GetUser(ctx context.Context, req *iam.GetUserRequest) (*iam.GetUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	user, err := i.repository.GetUser(req.UserId, req.Sub)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			appLogger.Infof("[GetUser]User not found: GetUserRequest=%+v", req)
			return &iam.GetUserResponse{}, nil
		}
		return nil, err
	}
	return &iam.GetUserResponse{User: convertUser(user)}, nil
}

func (i *iamService) PutUser(ctx context.Context, req *iam.PutUserRequest) (*iam.PutUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := i.repository.GetUserBySub(req.User.Sub)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var userID uint32
	if !noRecord {
		userID = savedData.UserID
	}
	u := &model.User{
		UserID:    userID,
		Sub:       req.User.Sub,
		Name:      req.User.Name,
		Activated: req.User.Activated,
	}

	// upsert
	registerdData, err := i.repository.PutUser(u)
	if err != nil {
		return nil, err
	}
	return &iam.PutUserResponse{User: convertUser(registerdData)}, nil
}

func convertUser(u *model.User) *iam.User {
	return &iam.User{
		UserId:    u.UserID,
		Sub:       u.Sub,
		Name:      u.Name,
		Activated: u.Activated,
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
}
