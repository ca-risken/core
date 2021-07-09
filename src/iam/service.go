package main

import (
	"context"
	"errors"
	"regexp"

	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

type iamService struct {
	repository iamRepository
}

func newIAMService() iam.IAMServiceServer {
	return &iamService{
		repository: newIAMRepository(),
	}
}

func (i *iamService) IsAuthorized(ctx context.Context, req *iam.IsAuthorizedRequest) (*iam.IsAuthorizedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	policies, err := i.repository.GetUserPolicy(req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

func (i *iamService) IsAdmin(ctx context.Context, req *iam.IsAdminRequest) (*iam.IsAdminResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	policy, err := i.repository.GetAdminPolicy(req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &iam.IsAdminResponse{Ok: false}, nil
		}
		return nil, err
	}
	appLogger.Debugf("user(%d) is admin, policy_id: %d", req.UserId, policy.PolicyID)
	return &iam.IsAdminResponse{Ok: true}, nil
}
