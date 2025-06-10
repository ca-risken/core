package organization_iam

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/iam"
	"github.com/ca-risken/core/proto/organization_iam"
	"gorm.io/gorm"
)

var (
	actionNamePattern = regexp.MustCompile(`^(\w|-)+/(\w|-)+$`)
)

func (i *OrganizationIAMService) IsAuthorizedOrganization(ctx context.Context, req *organization_iam.IsAuthorizedOrganizationRequest) (*organization_iam.IsAuthorizedOrganizationResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	isAdmin, err := i.iamClient.IsAdmin(ctx, &iam.IsAdminRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	if isAdmin.Ok {
		i.logger.Infof(ctx, "Authorized admin user action, request=%+v", req)
		return &organization_iam.IsAuthorizedOrganizationResponse{Ok: true}, nil
	}
	if !actionNamePattern.MatchString(req.ActionName) {
		return nil, fmt.Errorf("invalid action name, pattern=%s, action_name=%s", actionNamePattern, req.ActionName)
	}
	policies, err := i.repository.GetOrganizationPolicyByUserID(ctx, req.UserId, req.OrganizationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.IsAuthorizedOrganizationResponse{Ok: false}, nil
		}
		return nil, err
	}
	isAuthorized, err := isAuthorizedByOrganizationPolicy(req.ActionName, policies)
	if err != nil {
		return &organization_iam.IsAuthorizedOrganizationResponse{Ok: false}, err
	}
	if isAuthorized {
		i.logger.Infof(ctx, "Authorized organization action, request=%+v", req)
	}
	return &organization_iam.IsAuthorizedOrganizationResponse{Ok: isAuthorized}, nil
}

func isAuthorizedByOrganizationPolicy(action string, policies *[]model.OrganizationPolicy) (bool, error) {
	for _, p := range *policies {
		actionPtn, err := regexp.Compile(p.ActionPtn)
		if err != nil {
			return false, err
		}
		if actionPtn.MatchString(action) {
			return true, nil
		}
	}
	return false, nil
}
