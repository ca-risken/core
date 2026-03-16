package org_iam

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/org_iam"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (i *OrgIAMService) ListOrgPolicy(ctx context.Context, req *org_iam.ListOrgPolicyRequest) (*org_iam.ListOrgPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := i.repository.ListOrgPolicy(ctx, req.OrganizationId, req.Name, req.RoleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_iam.ListOrgPolicyResponse{}, nil
		}
		return nil, err
	}
	ids := []uint32{}
	for _, p := range list {
		ids = append(ids, p.PolicyID)
	}
	return &org_iam.ListOrgPolicyResponse{PolicyId: ids}, nil
}

func (i *OrgIAMService) GetOrgPolicy(ctx context.Context, req *org_iam.GetOrgPolicyRequest) (*org_iam.GetOrgPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	p, err := i.repository.GetOrgPolicy(ctx, req.OrganizationId, req.PolicyId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &org_iam.GetOrgPolicyResponse{}, nil
		}
		return nil, err
	}
	return &org_iam.GetOrgPolicyResponse{Policy: convertOrgPolicy(p)}, nil
}

func convertOrgPolicy(p *model.OrganizationPolicy) *org_iam.OrgPolicy {
	return &org_iam.OrgPolicy{
		PolicyId:       p.PolicyID,
		Name:           p.Name,
		OrganizationId: p.OrganizationID,
		ActionPtn:      p.ActionPtn,
		CreatedAt:      p.CreatedAt.Unix(),
		UpdatedAt:      p.UpdatedAt.Unix(),
	}
}

func (i *OrgIAMService) PutOrgPolicy(ctx context.Context, req *org_iam.PutOrgPolicyRequest) (*org_iam.PutOrgPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if _, err := regexp.Compile(req.ActionPtn); err != nil {
		return nil, fmt.Errorf("could not regexp complie, pattern=%s, err=%+v", req.ActionPtn, err)
	}
	savedData, err := i.repository.GetOrgPolicyByName(ctx, req.OrganizationId, req.Name)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, err
	}
	var policyID uint32
	if !noRecord {
		policyID = savedData.PolicyID
	}
	p := &model.OrganizationPolicy{
		PolicyID:       policyID,
		Name:           req.Name,
		OrganizationID: req.OrganizationId,
		ActionPtn:      req.ActionPtn,
	}
	registerdData, err := i.repository.PutOrgPolicy(ctx, p)
	if err != nil {
		return nil, err
	}
	return &org_iam.PutOrgPolicyResponse{Policy: convertOrgPolicy(registerdData)}, nil
}

func (i *OrgIAMService) DeleteOrgPolicy(ctx context.Context, req *org_iam.DeleteOrgPolicyRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := i.repository.DeleteOrgPolicy(ctx, req.OrganizationId, req.PolicyId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *OrgIAMService) AttachOrgPolicy(ctx context.Context, req *org_iam.AttachOrgPolicyRequest) (*org_iam.AttachOrgPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	policy, err := s.repository.AttachOrgPolicy(ctx, req.OrganizationId, req.PolicyId, req.RoleId)
	if err != nil {
		return nil, err
	}
	return &org_iam.AttachOrgPolicyResponse{Policy: convertOrgPolicy(policy)}, nil
}

func (s *OrgIAMService) DetachOrgPolicy(ctx context.Context, req *org_iam.DetachOrgPolicyRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.DetachOrgPolicy(ctx, req.OrganizationId, req.PolicyId, req.RoleId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
