package finding

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ca-risken/core/proto/finding"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

func (f *FindingService) GetAISummary(ctx context.Context, req *finding.GetAISummaryRequest) (*finding.GetAISummaryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if f.ai == nil {
		return nil, errors.New("unsupported AI service")
	}
	// Alert summaries read from master so the saved DB cache is visible immediately.
	data, err := f.repository.GetFinding(ctx, req.ProjectId, req.FindingId, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no finding: project_id=%d, finding_id=%d", req.ProjectId, req.FindingId)
		}
		return nil, err
	}
	recommend, err := f.repository.GetRecommend(ctx, req.ProjectId, req.FindingId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// UI summaries are always generated on demand and never read from DB cache.
	answer, err := f.ai.AskAISummaryFromFinding(ctx, data, recommend, req.Lang)
	if err != nil {
		return nil, fmt.Errorf("openai API error: err=%w", err)
	}
	return &finding.GetAISummaryResponse{Answer: answer}, nil
}

func (f *FindingService) GetAlertAISummary(ctx context.Context, req *finding.GetAlertAISummaryRequest) (*finding.GetAlertAISummaryResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if f.ai == nil {
		return nil, errors.New("unsupported AI service")
	}
	data, err := f.repository.GetFinding(ctx, req.ProjectId, req.FindingId, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no finding: project_id=%d, finding_id=%d", req.ProjectId, req.FindingId)
		}
		return nil, err
	}
	if data.AISummary != nil && *data.AISummary != "" {
		return &finding.GetAlertAISummaryResponse{AiSummary: *data.AISummary}, nil
	}

	recommend, err := f.repository.GetRecommend(ctx, req.ProjectId, req.FindingId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	answer, err := f.ai.AskAlertAISummaryFromFinding(ctx, data, recommend, req.Lang)
	if err != nil {
		return nil, fmt.Errorf("openai API error: err=%w", err)
	}

	now := time.Now()
	// finding.ai_summary is reserved for alert-summary cache only.
	if err := f.repository.UpdateFindingAISummary(ctx, req.ProjectId, req.FindingId, answer, now); err != nil {
		f.logger.Errorf(ctx, "Failed to save alert AI summary. project_id=%d finding_id=%d err=%v", req.ProjectId, req.FindingId, err)
	}
	return &finding.GetAlertAISummaryResponse{AiSummary: answer}, nil
}

func (f *FindingService) GetAISummaryStream(req *finding.GetAISummaryRequest, stream finding.FindingService_GetAISummaryStreamServer) error {
	if err := req.Validate(); err != nil {
		return err
	}
	if f.ai == nil {
		return errors.New("unsupported AI service")
	}
	ctx := context.Background()
	data, err := f.repository.GetFinding(ctx, req.ProjectId, req.FindingId, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no finding: project_id=%d, finding_id=%d", req.ProjectId, req.FindingId)
		}
		return err
	}
	recommend, err := f.repository.GetRecommend(ctx, req.ProjectId, req.FindingId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err := f.ai.AskAISummaryStreamFromFinding(ctx, data, recommend, req.Lang, stream); err != nil {
		return fmt.Errorf("openai API error: err=%w", err)
	}
	return nil
}

func (f *FindingService) UpdateFindingAISummary(ctx context.Context, req *finding.UpdateFindingAISummaryRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// This API is intended for alert-summary persistence only.
	if err := f.repository.UpdateFindingAISummary(ctx, req.ProjectId, req.FindingId, req.AiSummary, time.Unix(req.AiSummaryCreatedAt, 0)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "no finding: project_id=%d, finding_id=%d", req.ProjectId, req.FindingId)
		}
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
