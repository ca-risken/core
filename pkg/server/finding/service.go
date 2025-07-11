package finding

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/ai"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/finding"
)

const (
	defaultSortDirection string = "asc"
	defaultLimit         int32  = 200
)

var _ finding.FindingServiceServer = (*FindingService)(nil)

type FindingService struct {
	repository              db.FindingRepository
	logger                  logging.Logger
	ai                      ai.AIService
	excludeDeleteDataSource []string
}

var _ finding.FindingServiceServer = (*FindingService)(nil)

func NewFindingService(repository db.FindingRepository, openaiToken, chatGPTModel string, excludeDeleteDataSource []string, logger logging.Logger) *FindingService {
	return &FindingService{
		repository:              repository,
		logger:                  logger,
		ai:                      ai.NewAIClient(repository, openaiToken, chatGPTModel, logger),
		excludeDeleteDataSource: excludeDeleteDataSource,
	}
}

// TODO remove after fix response interface type change to int64 from uint32
func convertToUint32(v int64) uint32 {
	return uint32(v)
}
