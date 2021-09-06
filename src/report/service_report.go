package main

import (
	"context"
	"errors"
	"strings"

	"github.com/ca-risken/core/pkg/model"
	"github.com/ca-risken/core/proto/report"
	"github.com/golang/protobuf/ptypes/empty"
	"gorm.io/gorm"
)

/**
 * Report
 */

func (f *reportService) GetReportFinding(ctx context.Context, req *report.GetReportFindingRequest) (*report.GetReportFindingResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := f.repository.GetReportFinding(ctx, req.ProjectId, req.DataSource, req.FromDate, req.ToDate, req.Score)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &report.GetReportFindingResponse{}, nil
		}
		return nil, err
	}
	data := report.GetReportFindingResponse{}
	for _, d := range *list {
		data.ReportFinding = append(data.ReportFinding, convertReportFinding(&d))
	}
	return &data, nil
}

func (f *reportService) GetReportFindingAll(ctx context.Context, req *report.GetReportFindingAllRequest) (*report.GetReportFindingAllResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	list, err := f.repository.GetReportFindingAll(ctx, req.DataSource, req.FromDate, req.ToDate, req.Score)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &report.GetReportFindingAllResponse{}, nil
		}
		return nil, err
	}
	data := report.GetReportFindingAllResponse{}
	for _, d := range *list {
		data.ReportFinding = append(data.ReportFinding, convertReportFinding(&d))
	}
	return &data, nil
}

func (f *reportService) CollectReportFinding(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	err := f.repository.CollectReportFinding(ctx)
	if err != nil {
		appLogger.Errorf("Failed collectReportFinding. %v", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}

/**
 * Converter
 */

func convertReportFinding(f *model.ReportFinding) *report.ReportFinding {
	if f == nil {
		return &report.ReportFinding{}
	}
	return &report.ReportFinding{
		ReportFindingId: f.ReportFindingID,
		ReportDate:      f.ReportDate,
		ProjectId:       f.ProjectID,
		ProjectName:     f.ProjectName,
		Category:        strings.Split(f.DataSource, ":")[0],
		DataSource:      f.DataSource,
		Score:           f.Score,
		Count:           f.Count,
	}
}
