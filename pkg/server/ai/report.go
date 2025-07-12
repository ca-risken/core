package ai

import (
	"context"
	"errors"

	"github.com/ca-risken/core/proto/ai"
)

func (a *AIService) GenerateReport(ctx context.Context, req *ai.GenerateReportRequest) (*ai.GenerateReportResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if a.aiClient == nil {
		return nil, errors.New("unsupported AI service")
	}
	resp, err := a.aiClient.GenerateReport(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
