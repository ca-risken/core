package main

import (
	"context"
	"regexp"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/jinzhu/gorm"
	"github.com/vikyd/zero"
)

type iamService struct {
	repository iamRepoInterface
}

func newIAMService() iam.IAMServiceServer {
	return &iamService{
		repository: newIAMRepository(),
	}
}

func (i *iamService) GetUser(ctx context.Context, req *iam.GetUserRequest) (*iam.GetUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	user, err := i.repository.GetUser(req.UserId, req.Sub)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &iam.GetUserResponse{}, nil
		}
		return nil, err
	}
	return &iam.GetUserResponse{User: convertUser(user)}, nil
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

func (i *iamService) IsAuthorized(ctx context.Context, req *iam.IsAuthorizedRequest) (*iam.IsAuthorizedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	policies, err := i.repository.GetUserPoicy(req.UserId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &iam.IsAuthorizedResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized := false
	for _, p := range *policies {
		action, err := regexp.Compile(p.ActionPtn)
		if err != nil {
			return nil, err
		}
		resource, err := regexp.Compile(p.ResourcePtn)
		if err != nil {
			return nil, err
		}
		if !zero.IsZeroVal(p.ProjectID) && req.ProjectId != p.ProjectID {
			continue
		}
		if action.MatchString(req.ActionName) && resource.MatchString(req.ResourceName) {
			appLogger.Infof("Authorized user action, request=%+v", req)
			isAuthorized = true
			break
		}
	}
	return &iam.IsAuthorizedResponse{Ok: isAuthorized}, nil
}
