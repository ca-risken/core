package main

import (
	"context"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

func (f *findingService) GetPendFinding(ctx context.Context, req *finding.GetPendFindingRequest) (*finding.GetPendFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetPendFinding(req.ProjectId, req.FindingId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &finding.GetPendFindingResponse{}, nil
		}
		return nil, err
	}
	return &finding.GetPendFindingResponse{PendFinding: convertPendFinding(data)}, nil
}

func (f *findingService) PutPendFinding(ctx context.Context, req *finding.PutPendFindingRequest) (*finding.PutPendFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	registerd, err := f.repository.UpsertPendFinding(req.PendFinding.FindingId, req.PendFinding.ProjectId)
	if err != nil {
		return nil, err
	}
	return &finding.PutPendFindingResponse{PendFinding: convertPendFinding(registerd)}, nil
}

func (f *findingService) DeletePendFinding(ctx context.Context, req *finding.DeletePendFindingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := f.repository.DeletePendFinding(req.ProjectId, req.FindingId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertPendFinding(f *model.PendFinding) *finding.PendFinding {
	if f == nil {
		return &finding.PendFinding{}
	}
	return &finding.PendFinding{
		FindingId: f.FindingID,
		ProjectId: f.ProjectID,
		CreatedAt: f.CreatedAt.Unix(),
		UpdatedAt: f.UpdatedAt.Unix(),
	}
}