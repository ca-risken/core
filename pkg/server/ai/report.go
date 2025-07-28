package ai

import (
	"context"
	"errors"
	"time"

	"github.com/ca-risken/core/proto/ai"
	"github.com/ca-risken/core/proto/report"
)

func (a *AIService) GenerateReport(ctx context.Context, req *ai.GenerateReportRequest) (*ai.GenerateReportResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if a.aiClient == nil {
		return nil, errors.New("unsupported AI service")
	}

	// First, generate empty report
	r, err := a.putReport(ctx, 0, req.ProjectId, req.Name, "Markdown", "IN_PROGRESS", "")
	if err != nil {
		return nil, err
	}

	// Unsynchronized generation of report content
	go func() {
		// Generate report with AI
		genAICtx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
		defer cancel()
		content, err := a.aiClient.GenerateReport(genAICtx, req.ProjectId, req.Prompt)
		if err != nil {
			errCtx := context.Background()
			a.logger.Errorf(errCtx, "failed to generate report content: project_id=%d, report_id=%d, err=%v", req.ProjectId, r.Report.ReportId, err)
			if _, err := a.putReport(errCtx, r.Report.ReportId, req.ProjectId, req.Name, "Markdown", "ERROR", ""); err != nil {
				a.logger.Errorf(errCtx, "failed to update report content: project_id=%d, report_id=%d, err=%v", req.ProjectId, r.Report.ReportId, err)
			}
			return
		}
		_, err = a.putReport(genAICtx, r.Report.ReportId, req.ProjectId, req.Name, "Markdown", "OK", content)
		if err != nil {
			a.logger.Errorf(genAICtx, "failed to update report content: project_id=%d, report_id=%d, err=%v", req.ProjectId, r.Report.ReportId, err)
			return
		}
	}()

	// Early return
	return &ai.GenerateReportResponse{
		ReportId: r.Report.ReportId,
		Status:   r.Report.Status,
	}, nil
}

func (a *AIService) putReport(ctx context.Context, reportID, projectID uint32, name, reportType, status, content string) (*report.PutReportResponse, error) {
	var id uint32
	if reportID != 0 {
		id = reportID // update
	}

	return a.reportClient.PutReport(ctx, &report.PutReportRequest{
		ReportId:  id,
		ProjectId: projectID,
		Name:      name,
		Type:      reportType,
		Status:    status,
		Content:   content,
	})
}
