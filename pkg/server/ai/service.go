package ai

import (
	"github.com/ca-risken/common/pkg/logging"
	aiservice "github.com/ca-risken/core/pkg/ai"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/ai"
	"github.com/ca-risken/core/proto/report"
)

var _ ai.AIServiceServer = (*AIService)(nil)

type AIService struct {
	ai.UnimplementedAIServiceServer

	repository   db.AIRepository
	aiClient     aiservice.AIService
	reportClient report.ReportServiceClient
	logger       logging.Logger
}

func NewAIService(
	findingRepository db.FindingRepository,
	repository db.AIRepository,
	openaiToken string,
	chatGPTModel string,
	reasoningModel string,
	reportClient report.ReportServiceClient,
	logger logging.Logger,
) *AIService {
	return &AIService{
		repository:   repository,
		aiClient:     aiservice.NewAIClient(findingRepository, openaiToken, chatGPTModel, reasoningModel, logger),
		reportClient: reportClient,
		logger:       logger,
	}
}
