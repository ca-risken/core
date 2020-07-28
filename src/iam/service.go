package main

import (
	"context"
	"regexp"

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
