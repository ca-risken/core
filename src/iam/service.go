package main

import (
	"context"
	"errors"
	"regexp"

	"github.com/CyberAgent/mimosa-core/pkg/model"
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
	policies, err := i.repository.GetUserPolicy(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &iam.IsAuthorizedResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized, err := isAuthorizedByPolicy(req.ProjectId, req.ActionName, req.ResourceName, policies)
	if err != nil {
		return &iam.IsAuthorizedResponse{Ok: false}, err
	}
	if isAuthorized {
		appLogger.Infof("Authorized user action, request=%+v", req)
	}
	return &iam.IsAuthorizedResponse{Ok: isAuthorized}, nil
}

func (i *iamService) IsAuthorizedToken(ctx context.Context, req *iam.IsAuthorizedTokenRequest) (*iam.IsAuthorizedTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	existsMaintainer, err := i.repository.ExistsAccessTokenMaintainer(ctx, req.ProjectId, req.AccessTokenId)
	if err != nil {
		return nil, err
	}
	if !existsMaintainer {
		appLogger.Warnf("Unautorized the token that has no maintainers or expired in the project. project_id=%d, access_token_id=%d", req.ProjectId, req.AccessTokenId)
		return &iam.IsAuthorizedTokenResponse{Ok: false}, nil
	}
	policies, err := i.repository.GetTokenPolicy(ctx, req.AccessTokenId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &iam.IsAuthorizedTokenResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized, err := isAuthorizedByPolicy(req.ProjectId, req.ActionName, req.ResourceName, policies)
	if err != nil {
		return &iam.IsAuthorizedTokenResponse{Ok: false}, err
	}
	if isAuthorized {
		appLogger.Infof("Authorized access_token action, request=%+v", req)
	}
	return &iam.IsAuthorizedTokenResponse{Ok: isAuthorized}, nil
}

func isAuthorizedByPolicy(projectID uint32, action, resource string, policies *[]model.Policy) (bool, error) {
	for _, p := range *policies {
		actionPtn, err := regexp.Compile(p.ActionPtn)
		if err != nil {
			return false, err
		}
		resourcePtn, err := regexp.Compile(p.ResourcePtn)
		if err != nil {
			return false, err
		}
		if !zero.IsZeroVal(p.ProjectID) && projectID != p.ProjectID {
			continue
		}
		if actionPtn.MatchString(action) && resourcePtn.MatchString(resource) {
			return true, nil
		}
	}
	return false, nil

}

func (i *iamService) IsAdmin(ctx context.Context, req *iam.IsAdminRequest) (*iam.IsAdminResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	policy, err := i.repository.GetAdminPolicy(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &iam.IsAdminResponse{Ok: false}, nil
		}
		return nil, err
	}
	appLogger.Debugf("user(%d) is admin, policy_id: %d", req.UserId, policy.PolicyID)
	return &iam.IsAdminResponse{Ok: true}, nil
}
