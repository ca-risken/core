package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/src/finding/model"
	"github.com/vikyd/zero"
	"gorm.io/gorm"
)

func convertRecommend(findingID uint64, data *model.Recommend) *finding.Recommend {
	if data == nil {
		return &finding.Recommend{}
	}
	return &finding.Recommend{
		FindingId:      findingID,
		RecommendId:    data.RecommendID,
		DataSource:     data.DataSource,
		Type:           data.Type,
		Risk:           data.Risk,
		Recommendation: data.Recommendation,
		CreatedAt:      data.CreatedAt.Unix(),
		UpdatedAt:      data.CreatedAt.Unix(),
	}
}

func (f *findingService) GetRecommend(ctx context.Context, req *finding.GetRecommendRequest) (*finding.GetRecommendResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	data, err := f.repository.GetRecommend(ctx, req.ProjectId, req.FindingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &finding.GetRecommendResponse{}, nil
		}
		return nil, err
	}
	return &finding.GetRecommendResponse{Recommend: convertRecommend(req.FindingId, data)}, nil
}

func (f *findingService) PutRecommend(ctx context.Context, req *finding.PutRecommendRequest) (*finding.PutRecommendResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	// exists finding in the req project
	if f, err := f.repository.GetFinding(ctx, req.ProjectId, req.FindingId, true); err != nil || zero.IsZeroVal(f.FindingID) {
		appLogger.Warnf("Failed to get finding, project_id=%d, finding_id=%d, err=%+v", req.ProjectId, req.FindingId, err)
		return nil, err
	}
	registered, err := f.repository.UpsertRecommend(ctx, &model.Recommend{
		DataSource:     req.DataSource,
		Type:           req.Type,
		Risk:           req.Risk,
		Recommendation: req.Recommendation,
	})
	if err != nil {
		return nil, err
	}
	if _, err := f.repository.UpsertRecommendFinding(ctx, &model.RecommendFinding{
		FindingID:   req.FindingId,
		RecommendID: registered.RecommendID,
		ProjectID:   req.ProjectId,
	}); err != nil {
		return nil, err
	}
	return &finding.PutRecommendResponse{Recommend: convertRecommend(req.FindingId, registered)}, nil
}

func (f *findingService) getStoredRecommendID(ctx context.Context, dataSource, recommendType string) (*uint32, error) {
	storedData, err := f.repository.GetRecommendByDataSourceType(ctx, dataSource, recommendType)
	noRecord := errors.Is(err, gorm.ErrRecordNotFound)
	if err != nil && !noRecord {
		return nil, fmt.Errorf("Failed to GetRecommendByDataSourceType, data_source=%s, type=%s, err=%+v", dataSource, recommendType, err)
	}
	if !noRecord {
		return &storedData.RecommendID, nil
	}
	return nil, nil
}
