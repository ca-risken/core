package main

import (
	"context"
	"errors"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (f *findingService) GetPendFinding(ctx context.Context, req *finding.GetPendFindingRequest) (*finding.GetPendFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetPendFinding(req.ProjectId, req.FindingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
	_, err := f.repository.GetFinding(req.ProjectId, req.PendFinding.FindingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			appLogger.Warnf("RecordNotFound for PutPendFinding API, project_id=%d, finding=%d", req.ProjectId, req.PendFinding.FindingId)
		}
		return nil, err // DB error or RecordNotFound error
	}
	registerd, err := f.repository.UpsertPendFinding(req.PendFinding)
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
		Note:      f.Note,
		CreatedAt: f.CreatedAt.Unix(),
		UpdatedAt: f.UpdatedAt.Unix(),
	}
}
