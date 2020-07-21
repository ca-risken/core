package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/proto/iam"
	"github.com/golang/protobuf/ptypes/empty"
)

func (i *iamService) ListPolicy(ctx context.Context, req *iam.ListPolicyRequest) (*iam.ListPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *iamService) GetPolicy(ctx context.Context, req *iam.GetPolicyRequest) (*iam.GetPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *iamService) PutPolicy(ctx context.Context, req *iam.PutPolicyRequest) (*iam.PutPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *iamService) DeletePolicy(ctx context.Context, req *iam.DeletePolicyRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *iamService) AttachPolicy(ctx context.Context, req *iam.AttachPolicyRequest) (*iam.AttachPolicyResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *iamService) DetachPolicy(ctx context.Context, req *iam.DetachPolicyRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return nil, nil
}
