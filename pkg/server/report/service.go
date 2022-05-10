package report

import "github.com/ca-risken/core/pkg/db"

type ReportService struct {
	repository db.ReportRepository
}

func NewReportService(repository db.ReportRepository) *ReportService {
	return &ReportService{
		repository: repository,
	}
}
