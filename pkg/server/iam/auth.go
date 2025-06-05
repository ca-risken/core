package iam

import (
	"context"
	"errors"
	"regexp"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

// isUserAdmin checks if the user is admin
func (i *IAMService) isUserAdmin(ctx context.Context, userID uint32) (bool, error) {
	user, err := i.repository.GetUser(ctx, userID, "", "")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return user.IsAdmin, nil
}

func (i *IAMService) IsAuthorized(ctx context.Context, req *iam.IsAuthorizedRequest) (*iam.IsAuthorizedResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Check if user is admin first
	isAdmin, err := i.isUserAdmin(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	// If user is admin, always return true
	if isAdmin {
		i.logger.Infof(ctx, "Authorized admin user action, request=%+v", req)
		return &iam.IsAuthorizedResponse{Ok: true}, nil
	}

	// If user exists but is not admin, check policies
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
		i.logger.Infof(ctx, "Authorized user action, request=%+v", req)
	}
	return &iam.IsAuthorizedResponse{Ok: isAuthorized}, nil
}

func (i *IAMService) IsAuthorizedAdmin(ctx context.Context, req *iam.IsAuthorizedAdminRequest) (*iam.IsAuthorizedAdminResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	isAdmin, err := i.isUserAdmin(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &iam.IsAuthorizedAdminResponse{Ok: isAdmin}, nil
}

func (i *IAMService) IsAuthorizedToken(ctx context.Context, req *iam.IsAuthorizedTokenRequest) (*iam.IsAuthorizedTokenResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	existsMaintainer, err := i.repository.ExistsAccessTokenMaintainer(ctx, req.ProjectId, req.AccessTokenId)
	if err != nil {
		return nil, err
	}
	if !existsMaintainer {
		i.logger.Warnf(ctx, "Unautorized the token that has no maintainers or expired in the project. project_id=%d, access_token_id=%d", req.ProjectId, req.AccessTokenId)
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
		i.logger.Infof(ctx, "Authorized access_token action, request=%+v", req)
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

func (i *IAMService) IsAdmin(ctx context.Context, req *iam.IsAdminRequest) (*iam.IsAdminResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	isAdmin, err := i.isUserAdmin(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	i.logger.Debugf(ctx, "user(%d) is_admin: %t", req.UserId, isAdmin)
	return &iam.IsAdminResponse{Ok: isAdmin}, nil
}
