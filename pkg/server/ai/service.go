package ai

import (
	"github.com/ca-risken/common/pkg/logging"
	aiservice "github.com/ca-risken/core/pkg/ai"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/ai"
)

var _ ai.AIServiceServer = (*AIService)(nil)

type AIService struct {
	aiClient aiservice.AIService
	logger   logging.Logger
}

func NewAIService(
	repository db.FindingRepository,
	openaiToken string,
	chatGPTModel string,
	logger logging.Logger,
) *AIService {
	return &AIService{
		aiClient: aiservice.NewAIClient(repository, openaiToken, chatGPTModel, logger),
		logger:   logger,
	}
}
