package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/CyberAgent/mimosa-core/pkg/model"
	"github.com/CyberAgent/mimosa-core/proto/finding"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
)

func (f *findingService) ListFinding(ctx context.Context, req *finding.ListFindingRequest) (*finding.ListFindingResponse, error) {
	resp := finding.ListFindingResponse{}
	if err := req.Validate(); err != nil {
		return &resp, err
	}
	list, err := f.repository.ListFinding(convertListFindingRequest(req))
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &resp, nil
		}
		return &resp, err
	}

	// TODO authz
	var ids []uint64
	for _, id := range *list {
		ids = append(ids, uint64(id.FindingID))
	}
	resp.FindingId = ids
	return &resp, nil
}

func convertListFindingRequest(req *finding.ListFindingRequest) *finding.ListFindingRequest {
	converted := finding.ListFindingRequest{
		ProjectId:    req.ProjectId,
		ResourceName: req.ResourceName,
		DataSource:   req.DataSource,
		FromScore:    req.FromScore,
		ToScore:      req.ToScore,
		FromAt:       req.FromAt,
		ToAt:         req.ToAt,
	}
	if converted.ToScore == 0 {
		converted.ToScore = 1.0
	}
	if converted.ToAt == 0 {
		converted.ToAt = time.Now().Unix()
	}
	return &converted
}

func (f *findingService) GetFinding(ctx context.Context, req *finding.GetFindingRequest) (*finding.GetFindingResponse, error) {
	resp := finding.GetFindingResponse{}
	if err := req.Validate(); err != nil {
		return &resp, err
	}

	// TODO authz
	data, err := f.repository.GetFinding(req.FindingId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &resp, nil
		}
		return &resp, err
	}
	resp.Finding = convertFinding(data)
	return &resp, nil
}

func (f *findingService) PutFinding(ctx context.Context, req *finding.PutFindingRequest) (*finding.PutFindingResponse, error) {
	resp := finding.PutFindingResponse{}
	if err := req.Finding.Validate(); err != nil {
		return &resp, err
	}

	savedData, err := f.repository.GetFindingByDataSource(
		req.Finding.ProjectId, req.Finding.DataSource, req.Finding.DataSourceId)
	noRecord := gorm.IsRecordNotFoundError(err)
	if err != nil && !noRecord {
		return &resp, err
	}

	var findingID uint64
	if !noRecord {
		if req.Finding.FindingId != 0 && req.Finding.FindingId != savedData.FindingID {
			return &resp, fmt.Errorf("Invalid finding_id, want=%d, got=%d", savedData.FindingID, req.Finding.FindingId)
		}
		findingID = savedData.FindingID
	}

	// TODO: Authz
	// Fiding upsert
	registerdData, err := f.repository.UpsertFinding(
		&model.Finding{
			FindingID:     findingID,
			Description:   req.Finding.Description,
			DataSource:    req.Finding.DataSource,
			DataSourceID:  req.Finding.DataSourceId,
			ResourceName:  req.Finding.ResourceName,
			ProjectID:     req.Finding.ProjectId,
			OriginalScore: req.Finding.OriginalScore,
			Score:         calculateScore(req.Finding.OriginalScore, req.Finding.OriginalMaxScore),
			Data:          req.Finding.Data,
		})
	if err != nil {
		return &resp, err
	}

	// Resource
	_, err = f.PutResource(ctx, &finding.PutResourceRequest{
		Resource: &finding.ResourceForUpsert{
			ResourceName: req.Finding.ResourceName,
			ProjectId:    req.Finding.ProjectId,
		},
	})
	if err != nil {
		return &resp, err
	}
	resp.Finding = convertFinding(registerdData)
	return &resp, nil
}

func (f *findingService) DeleteFinding(ctx context.Context, req *finding.DeleteFindingRequest) (*empty.Empty, error) {
	resp := empty.Empty{}
	if err := req.Validate(); err != nil {
		return &resp, err
	}

	// TODO authz
	err := f.repository.DeleteFinding(req.FindingId)
	if err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (f *findingService) ListFindingTag(ctx context.Context, req *finding.ListFindingTagRequest) (*finding.ListFindingTagResponse, error) {
	resp := finding.ListFindingTagResponse{}
	if err := req.Validate(); err != nil {
		return &resp, err
	}

	// TODO authz
	list, err := f.repository.ListFindingTag(req.FindingId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return &resp, nil
		}
		return &resp, err
	}
	var tags []*finding.FindingTag
	for _, tag := range *list {
		tags = append(tags, &finding.FindingTag{
			FindingTagId: tag.FindingTagID,
			FindingId:    tag.FindingID,
			TagKey:       tag.TagKey,
			TagValue:     tag.TagValue,
			CreatedAt:    tag.CreatedAt.Unix(),
			UpdatedAt:    tag.UpdatedAt.Unix(),
		})
	}
	resp.Tag = tags
	return &resp, nil
}

func (f *findingService) TagFinding(ctx context.Context, req *finding.TagFindingRequest) (*finding.TagFindingResponse, error) {
	resp := finding.TagFindingResponse{}
	if err := req.Validate(); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func (f *findingService) UntagFinding(ctx context.Context, req *finding.UntagFindingRequest) (*empty.Empty, error) {
	resp := empty.Empty{}
	if err := req.Validate(); err != nil {
		return &resp, err
	}
	return &resp, nil
}

func convertFinding(f *model.Finding) *finding.Finding {
	if f == nil {
		return &finding.Finding{}
	}
	return &finding.Finding{
		FindingId:     f.FindingID,
		Description:   f.Description,
		DataSource:    f.DataSource,
		DataSourceId:  f.DataSourceID,
		ResourceName:  f.ResourceName,
		ProjectId:     f.ProjectID,
		OriginalScore: f.OriginalScore,
		Score:         f.Score,
		Data:          f.Data,
		CreatedAt:     f.CreatedAt.Unix(),
		UpdatedAt:     f.UpdatedAt.Unix(),
	}
}

func calculateScore(score, maxScore float32) float32 {
	return float32(math.Round(float64(score/maxScore*100)) / 100)
}
