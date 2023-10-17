package finding

import (
	"context"
	"errors"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

func (f *FindingService) GetPendFinding(ctx context.Context, req *finding.GetPendFindingRequest) (*finding.GetPendFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetPendFinding(ctx, req.ProjectId, req.FindingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &finding.GetPendFindingResponse{}, nil
		}
		return nil, err
	}
	return &finding.GetPendFindingResponse{PendFinding: convertPendFinding(data)}, nil
}

func (f *FindingService) PutPendFinding(ctx context.Context, req *finding.PutPendFindingRequest) (*finding.PutPendFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	_, err := f.repository.GetFinding(ctx, req.ProjectId, req.PendFinding.FindingId, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			f.logger.Warnf(ctx, "RecordNotFound for PutPendFinding API, project_id=%d, finding=%d", req.ProjectId, req.PendFinding.FindingId)
		}
		return nil, err // DB error or RecordNotFound error
	}
	registerd, err := f.repository.UpsertPendFinding(ctx, req.PendFinding)
	if err != nil {
		return nil, err
	}
	return &finding.PutPendFindingResponse{PendFinding: convertPendFinding(registerd)}, nil
}

func (f *FindingService) DeletePendFinding(ctx context.Context, req *finding.DeletePendFindingRequest) (*empty.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if err := f.repository.DeletePendFinding(ctx, req.ProjectId, req.FindingId); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func convertPendFinding(f *model.PendFinding) *finding.PendFinding {
	if f == nil {
		return &finding.PendFinding{}
	}
	converted := &finding.PendFinding{
		FindingId:     f.FindingID,
		ProjectId:     f.ProjectID,
		Note:          f.Note,
		FalsePositive: f.FalsePositive,
		CreatedAt:     f.CreatedAt.Unix(),
		UpdatedAt:     f.UpdatedAt.Unix(),
	}
	if !f.ExpiredAt.IsZero() {
		converted.ExpiredAt = f.ExpiredAt.Unix()
	}
	return converted
}
