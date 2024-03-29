package finding

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/proto/finding"
	"gorm.io/gorm"
)

func (f *FindingService) GetAISummary(ctx context.Context, req *finding.GetAISummaryRequest) (*finding.GetAISummaryResponse, error) {
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
	recommend, err := f.repository.GetRecommend(ctx, req.ProjectId, req.FindingId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	answer, err := f.ai.AskAISummaryFromFinding(ctx, data, recommend, req.Lang)
	if err != nil {
		return nil, fmt.Errorf("openai API error: err=%w", err)
	}
	return &finding.GetAISummaryResponse{Answer: answer}, nil
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
