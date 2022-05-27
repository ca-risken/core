package report

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
)

type ReportService struct {
	repository db.ReportRepository
	logger     logging.Logger
}

func NewReportService(repository db.ReportRepository, logger logging.Logger) *ReportService {
	return &ReportService{
		repository: repository,
		logger:     logger,
	}
}
