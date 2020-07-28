package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

func (i *iamService) ListPolicy(ctx context.Context, req *iam.ListPolicyRequest) (*iam.ListPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListPolicy(req.ProjectId, req.Name)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &iam.ListPolicyResponse{}, nil
		}
		return nil, err
	}
	ids := []uint32{}
	for _, p := range *list {
		ids = append(ids, p.PolicyID)
	}
	return &iam.ListPolicyResponse{PolicyId: ids}, nil
}

func (i *iamService) GetPolicy(ctx context.Context, req *iam.GetPolicyRequest) (*iam.GetPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	p, err := i.repository.GetPolicy(req.ProjectId, req.PolicyId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &iam.GetPolicyResponse{}, nil
		}
		return nil, err
	}
	return &iam.GetPolicyResponse{Policy: convertPolicy(p)}, nil
}

func convertPolicy(p *model.Policy) *iam.Policy {
	return &iam.Policy{
		PolicyId:    p.PolicyID,
		Name:        p.Name,
		ProjectId:   p.ProjectID,
		ActionPtn:   p.ActionPtn,
		ResourcePtn: p.ResourcePtn,
		CreatedAt:   p.CreatedAt.Unix(),
		UpdatedAt:   p.UpdatedAt.Unix(),
	}
}

func (i *iamService) PutPolicy(ctx context.Context, req *iam.PutPolicyRequest) (*iam.PutPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	savedData, err := i.repository.GetPolicyByName(req.Policy.ProjectId, req.Policy.Name)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return nil, err
	}

	// PKが登録済みの場合は取得した値をセット。未登録はゼロ値のママでAutoIncrementさせる（更新の都度、無駄にAutoIncrementさせないように）
	var policyID uint32
	if !noRecord {
		policyID = savedData.PolicyID
	}
	p := &model.Policy{
		PolicyID:    policyID,
		Name:        req.Policy.Name,
		ProjectID:   req.Policy.ProjectId,
		ActionPtn:   req.Policy.ActionPtn,
		ResourcePtn: req.Policy.ResourcePtn,
	}

	// upsert
	registerdData, err := i.repository.PutPolicy(p)
	if err != nil {
		return nil, err
	}
	return &iam.PutPolicyResponse{Policy: convertPolicy(registerdData)}, nil
}

func (i *iamService) DeletePolicy(ctx context.Context, req *iam.DeletePolicyRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DeletePolicy(req.ProjectId, req.PolicyId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (i *iamService) AttachPolicy(ctx context.Context, req *iam.AttachPolicyRequest) (*iam.AttachPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	rp, err := i.repository.AttachPolicy(req.ProjectId, req.RoleId, req.PolicyId)
	if err != nil {
		return nil, err
	}
	return &iam.AttachPolicyResponse{RolePolicy: convertRolePolicy(rp)}, nil
}

func convertRolePolicy(rp *model.RolePolicy) *iam.RolePolicy {
	return &iam.RolePolicy{
		RoleId:    rp.RoleID,
		PolicyId:  rp.PolicyID,
		ProjectId: rp.ProjectID,
		CreatedAt: rp.CreatedAt.Unix(),
		UpdatedAt: rp.UpdatedAt.Unix(),
	}
}

func (i *iamService) DetachPolicy(ctx context.Context, req *iam.DetachPolicyRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DetachPolicy(req.ProjectId, req.RoleId, req.PolicyId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
