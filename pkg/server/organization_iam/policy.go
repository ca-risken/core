package organization_iam

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/organization_iam"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (i *OrganizationIAMService) ListOrganizationPolicy(ctx context.Context, req *organization_iam.ListOrganizationPolicyRequest) (*organization_iam.ListOrganizationPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListOrganizationPolicy(ctx, req.OrganizationId, req.Name, req.RoleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.ListOrganizationPolicyResponse{}, nil
		}
		return nil, err
	}
	ids := []uint32{}
	for _, p := range list {
		ids = append(ids, p.PolicyID)
	}
	return &organization_iam.ListOrganizationPolicyResponse{PolicyId: ids}, nil
}

func (i *OrganizationIAMService) GetOrganizationPolicy(ctx context.Context, req *organization_iam.GetOrganizationPolicyRequest) (*organization_iam.GetOrganizationPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	p, err := i.repository.GetOrganizationPolicy(ctx, req.OrganizationId, req.PolicyId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &organization_iam.GetOrganizationPolicyResponse{}, nil
		}
		return nil, err
	}
	return &organization_iam.GetOrganizationPolicyResponse{Policy: convertOrganizationPolicy(p)}, nil
}

func convertOrganizationPolicy(p *model.OrganizationPolicy) *organization_iam.OrganizationPolicy {
	return &organization_iam.OrganizationPolicy{
		PolicyId:       p.PolicyID,
		Name:           p.Name,
		OrganizationId: p.OrganizationID,
		ActionPtn:      p.ActionPtn,
		CreatedAt:      p.CreatedAt.Unix(),
		UpdatedAt:      p.UpdatedAt.Unix(),
	}
}

func (i *OrganizationIAMService) PutOrganizationPolicy(ctx context.Context, req *organization_iam.PutOrganizationPolicyRequest) (*organization_iam.PutOrganizationPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.OrganizationId != req.Policy.OrganizationId {
		return nil, fmt.Errorf("unexpected organization_id: organization_id=%d, policy.organization_id=%d", req.OrganizationId, req.Policy.OrganizationId)
	}
	if _, err := regexp.Compile(req.Policy.ActionPtn); err != nil {
		return nil, fmt.Errorf("could not regexp complie, pattern=%s, err=%+v", req.Policy.ActionPtn, err)
	}
	savedData, err := i.repository.GetOrganizationPolicyByName(ctx, req.Policy.OrganizationId, req.Policy.Name)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, err
	}

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var policyID uint32
	if !noRecord {
		policyID = savedData.PolicyID
	}
	p := &model.OrganizationPolicy{
		PolicyID:       policyID,
		Name:           req.Policy.Name,
		OrganizationID: req.Policy.OrganizationId,
		ActionPtn:      req.Policy.ActionPtn,
	}

	// upsert
	registerdData, err := i.repository.PutOrganizationPolicy(ctx, p)
	if err != nil {
		return nil, err
	}
	return &organization_iam.PutOrganizationPolicyResponse{Policy: convertOrganizationPolicy(registerdData)}, nil
}

func (i *OrganizationIAMService) DeleteOrganizationPolicy(ctx context.Context, req *organization_iam.DeleteOrganizationPolicyRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DeleteOrganizationPolicy(ctx, req.OrganizationId, req.PolicyId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
