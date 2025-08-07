package report

import (
	"context"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/report"
)

func convertReport(r *model.Report) *report.Report {
	return &report.Report{
		ReportId:  r.ReportID,
		ProjectId: r.ProjectID,
		Name:      r.Name,
		Type:      r.Type,
		Status:    r.Status,
		Content:   r.Content,
		CreatedAt: r.CreatedAt.Unix(),
		UpdatedAt: r.UpdatedAt.Unix(),
	}
}

func (f *ReportService) GetReport(ctx context.Context, req *report.GetReportRequest) (*report.GetReportResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	r, err := f.repository.GetReport(ctx, req.ProjectId, req.ReportId)
	if err != nil {
		return nil, err
	}
	return &report.GetReportResponse{Report: convertReport(r)}, nil
}

func (f *ReportService) ListReport(ctx context.Context, req *report.ListReportRequest) (*report.ListReportResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := f.repository.ListReport(ctx, req.ProjectId)
	if err != nil {
		return nil, err
	}
	data := report.ListReportResponse{}
	for _, r := range *list {
		data.Report = append(data.Report, convertReport(&r))
	}
	return &data, nil
}

func (f *ReportService) PutReport(ctx context.Context, req *report.PutReportRequest) (*report.PutReportResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Get existing report
	if req.ReportId != 0 {
		_, err := f.repository.GetReport(ctx, req.ProjectId, req.ReportId)
		if err != nil {
			return nil, err
		}
	}

	var reportID uint32
	if req.ReportId != 0 {
		reportID = req.ReportId // if exists, use existing report_id(update)
	}

	// Upsert report
	r, err := f.repository.PutReport(ctx, &model.Report{
		ReportID:  reportID,
		ProjectID: req.ProjectId,
		Name:      req.Name,
		Type:      req.Type,
		Status:    req.Status,
		Content:   req.Content,
	})
	if err != nil {
		return nil, err
	}
	return &report.PutReportResponse{Report: convertReport(r)}, nil
}
