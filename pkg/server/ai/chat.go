package ai

import (
	"context"
	"errors"

	"github.com/ca-risken/core/proto/ai"
)

func (a *AIService) ChatAI(ctx context.Context, req *ai.ChatAIRequest) (*ai.ChatAIResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if a.aiClient == nil {
		return nil, errors.New("unsupported AI service")
	}
	resp, err := a.aiClient.ChatAI(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
