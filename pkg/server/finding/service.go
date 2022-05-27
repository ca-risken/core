package finding

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
)

const (
	defaultSortDirection string = "asc"
	defaultLimit         int32  = 200
)

type FindingService struct {
	repository db.FindingRepository
	logger     logging.Logger
}

func NewFindingService(repository db.FindingRepository, logger logging.Logger) *FindingService {
	return &FindingService{
		repository: repository,
		logger:     logger,
	}
}

// TODO remove after fix response interface type change to int64 from uint32
func convertToUint32(v int64) uint32 {
	return uint32(v)
}
