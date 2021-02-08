package main

import (
	"github.com/CyberAgent/mimosa-core/proto/report"
)

type reportService struct {
	repository reportRepository
}

func newReportService() report.ReportServiceServer {
	return &reportService{
		repository: newReportRepository(),
	}
}
