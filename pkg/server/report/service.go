package report

import (
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/pkg/db"
	"github.com/ca-risken/core/proto/report"
)

var _ report.ReportServiceServer = (*ReportService)(nil)

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
